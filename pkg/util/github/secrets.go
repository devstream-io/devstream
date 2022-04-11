package github

import (
	"context"
	crypto_rand "crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/google/go-github/v42/github"
	"golang.org/x/crypto/nacl/box"
)

// AddRepoSecret adds a secret to a GitHub repo.
func (c *Client) AddRepoSecret(secretKey, secretValue string) error {
	var owner = c.Owner
	if c.Org != "" {
		owner = c.Org
	}

	// The transmission of the secret value to GitHub using the api requires the secret value to be encrypted
	// with the public key of the repo.
	// First, the public key of the repo is retrieved.
	ctx := context.Background()
	publicKey, _, err := client.Actions.GetRepoPublicKey(ctx, owner, c.Repo)
	if err != nil {
		return err
	}

	// Second, we encrypt the secret to get a github.EncodedSecret
	encryptedSecret, err := encryptSecretWithPublicKey(publicKey, secretKey, secretValue)
	if err != nil {
		return err
	}

	// Finally, the github.EncodedSecret is passed into the GitHub client.Actions.CreateOrUpdateRepoSecret method to
	// create the secret in the GitHub repo
	if _, err := client.Actions.CreateOrUpdateRepoSecret(ctx, owner, c.Repo, encryptedSecret); err != nil {
		return fmt.Errorf("github Actions.CreateOrUpdateRepoSecret returned error: %v", err)
	}

	return nil
}

func encryptSecretWithPublicKey(publicKey *github.PublicKey, secretName string, secretValue string) (*github.EncryptedSecret, error) {
	// The public key comes base64 encoded, so it must be decoded prior to use.
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey.GetKey())
	if err != nil {
		return nil, fmt.Errorf("base64.StdEncoding.DecodeString was unable to decode public key: %v", err)
	}

	// The decode key is converted into a fixed size byte array.
	var boxKey [32]byte

	// The secret value is converted into a slice of bytes.
	copy(boxKey[:], decodedPublicKey)
	secretBytes := []byte(secretValue)

	// The secret is encrypted with box.SealAnonymous using the repo's decoded public key.
	encryptedBytes, err := box.SealAnonymous([]byte{}, secretBytes, &boxKey, crypto_rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("box.SealAnonymous failed with error %w", err)
	}

	// The encrypted secret is encoded as a base64 string to be used in a github.EncodedSecret type.
	encryptedString := base64.StdEncoding.EncodeToString(encryptedBytes)
	keyID := publicKey.GetKeyID()
	encryptedSecret := &github.EncryptedSecret{
		Name:           secretName,
		KeyID:          keyID,
		EncryptedValue: encryptedString,
	}
	return encryptedSecret, nil
}

// DeleteRepoSecret deletes a secret in a GitHub repo.
func (c *Client) DeleteRepoSecret(secretKey string) error {
	var owner = c.Owner
	if c.Org != "" {
		owner = c.Org
	}

	ctx := context.Background()
	response, err := client.Actions.DeleteRepoSecret(ctx, owner, c.Repo, secretKey)
	if err != nil {
		if response.StatusCode == 404 {
			return nil
		}
		return fmt.Errorf("github Actions.DeleteRepoSecret returned error: %v", err)
	}

	return nil
}

// RepoSecretExists detects if a secret exists in a GitHub repo.
func (c *Client) RepoSecretExists(secretKey string) (bool, error) {
	var owner = c.Owner
	if c.Org != "" {
		owner = c.Org
	}

	ctx := context.Background()
	_, response, err := client.Actions.GetRepoSecret(ctx, owner, c.Repo, secretKey)
	if err != nil {
		if response.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("github Actions.GetRepoSecret returned error: %v", err)
	}
	return true, nil
}

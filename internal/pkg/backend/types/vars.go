package types

type Type string

func (t Type) String() string {
	return string(t)
}

const (
	Local   Type = "local"
	S3      Type = "s3"
	K8s     Type = "k8s"
	K8sAlis Type = "kubernetes"
)

curl -sL https://www.cloudflare.com/ips-v4 | grep -v '^$' >ips-v4
cat ips-v4
if git status --porcelain | grep -q "^ M";then
  echo "ips-v4 changed"
  echo "CHANGED=1" >> "$GITHUB_ENV"
  git config user.name 'github-actions[bot]'
  git config user.email 'github-actions[bot]@users.noreply.github.com'
  git add ips-v4
  git commit -m "chore: ips-v4 changed"
  git push
else
  echo "ips-v4 not changed"
fi

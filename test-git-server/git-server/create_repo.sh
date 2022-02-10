#!/bin/sh

if [[ ! -d "test-repo" ]]; then
  echo "test-repo doesn't exists already!"
  mkdir test-repo
fi

cd test-repo

git config --global user.email "test@docplanner.com"
git config --global user.name "test-user"
git config --global init.defaultBranch develop
git init --shared=true

if [[ ! -d "example-app" ]]; then
  echo "example-app doesn't exists already!"
  mkdir example-app
fi

cat <<EOF > example-app/values.yaml
image:
  tag: 1.0.0
EOF

git add .
git commit -m "my first commit"

cd ..
git clone --bare test-repo test-repo.git

cp -R test-repo.git /git-server/repos

chown -R git:git /git-server/repos/

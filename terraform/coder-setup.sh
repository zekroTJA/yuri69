sudo apt clean

# Install Golang
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update -y
sudo apt install -y golang-go

# Install nodejs & NPM
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

# Install Corepack / Yarn
sudo corepack enable
sudo npm i -g corepack

# Install ffmpeg
sudo apt install -y ffmpeg

# Install Taskfile
sudo env GOBIN=/bin go install github.com/go-task/task/v3/cmd/task@latest
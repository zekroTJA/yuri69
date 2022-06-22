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

#!/bin/bash

COMPOSE_URL="https://raw.githubusercontent.com/zekroTJA/yuri69/master/docker-compose.yml"
LAVALINK_URL="https://raw.githubusercontent.com/zekroTJA/yuri69/master/config/lavalink/application.yml"

check_installed() {
  program=$1

  which "$program" > /dev/null 2>&1 || {
    echo "error: $program is not intstalled but required to execute this script."
    echo "       Please install $program and re-run the script."
    exit 1
  }
}

substitute() {
  file=$1
  key=$2
  value=$3

  sed -i "s/{{$key}}/${value//\//\\\/}/" "$file"
}

read_yn() {
  read -r yn
  case "$yn" in
    "y"|"Y"|"yes") return 0 ;;
    *) return 1 ;;
  esac
}


check_installed docker
check_installed curl


echo -e "To set up yuri69, we need to know some parameters.\n"

echo -e "First of all, make sure you have created a Discord bot application. If not, you can do this here:"
echo -e "https://discord.com/developers/applications\n"

echo "Please enter your bot token."
printf "> "
read -rs discord_token

echo -e "\n\nNow, please enter the bot's client ID."
printf "> "
read -r client_id

echo -e "\nNow we need the client secret (not the bot's token!)."
printf "> "
read -rs client_secret

echo -e "\n\nPlease tell us your Discord user ID to set you as the bot's owner."
printf "> "
read -r owner_id

echo -e "\nPlease enter the domain on which the bot will be accessible."
echo -e "(i.e. 'yuri.zekro.de')"
printf "(leave empty to use 'localhost') > "
read -r domain
[ -z "$domain" ] && domain="localhost"

echo -e "\nDo you want to expose the service via HTTPS?"
printf "(yN) > "
if read_yn; then
  echo -e "\nTo issue TLS certificates via ACME, we need your e-mail address."
  echo -e "This is only used as contact point for the ACME process."
  printf "> "
  read -r email

  http_entrypoint="https"
  enable_tls="true"
else
  http_entrypoint="http"
  enable_tls="false"
  email="hello@example.com"
fi

echo -e "\nThank you! Things are getting set up now ...\n"

set -ex

curl -Lo docker-compose.yml "$COMPOSE_URL"
substitute docker-compose.yml "ACME_EMAIL" "$email"
substitute docker-compose.yml "DISCORD_TOKEN" "$discord_token"
substitute docker-compose.yml "DISCORD_OWNER_ID" "$owner_id"
substitute docker-compose.yml "PUBLIC_ADDRESS" "$http_entrypoint://$domain"
substitute docker-compose.yml "DISCORD_CLIENT_ID" "$client_id"
substitute docker-compose.yml "DISCORD_CLIENT_SECRET" "$client_secret"
substitute docker-compose.yml "HTTP_ENTRYPOINT" "$http_entrypoint"
substitute docker-compose.yml "ENABLE_TLS" "$enable_tls"
substitute docker-compose.yml "PUBLIC_DOMAIN" "$domain"

mkdir -p config/lavalink
curl -Lo config/lavalink/application.yml "$LAVALINK_URL"

set +x

echo -e "\nEverything is set up and ready to start!"
echo -e "Do you want to start the service stack now?"
printf "(yN) > "
if read_yn; then
  docker compose up -d
  docker compose ps
else
  echo -e "\nYou can start the stack at any time by using the following command:"
  echo "docker compose up -d"
fi

echo -e "\nTo be able to log in via the web interface, you need to add the following URL to the OAuth2 Redirects section in your bot's application settings."
echo "(see: https://discord.com/developers/applications/$client_id/oauth2/general)"
echo -e "$http_entrypoint://$domain/api/v1/auth/oauth2/discord/callback"

echo -e "\nAfter that, you can log in to the Yuri69 web interface here:"
echo "$http_entrypoint://$domain"

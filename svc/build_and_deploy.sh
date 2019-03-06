set -e

GOOS=linux GOARCH=386 go build

gcloud compute scp svc devoted:~
gcloud compute ssh devoted --command="sudo mv svc /usr/sbin/devoted"

gcloud compute ssh devoted --command="sudo systemctl restart devoted.service"
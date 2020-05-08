rm -rf /sheck

while true; do
    gsutil -m rsync -rdu /sheck gs://chipmunk-storage/
    sleep 1
done
import subprocess
"""
 This program make 1000 http GET requests to the
 application container and ensures sure that the number
 returned by the container's counter matches the
 true count
"""
def main():
    counter = 0
    docker_ip = subprocess.getoutput(["docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' application"])
    print("Docker ip address: ", docker_ip)

    for i in range(1000):
        counter += 1
        current_num = int(subprocess.getoutput(["curl -s "+docker_ip+":8080"]))
        print(current_num, flush=True)
        if current_num != counter:
            print("FAILED", flush=True)
            exit(-1)

    print("TEST PASSED")

if __name__ == "__main__":
    main()

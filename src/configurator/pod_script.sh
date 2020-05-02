echo "starting configurator [pod]"

if [ -f /host/i_was_here ] && [ -f /host/bin/criu ]; then
    echo " -> goodbye forever [pod]"
    ./kubectl taint node --overwrite $(hostname) configured=true:NoSchedule

    export pod_name=$(./kubectl get pods --template '{{range .items}}{{.metadata.name}}{{"\n"}}{{end}}' --field-selector spec.nodeName=$(hostname) | grep node-configurator)

    echo $(./kubectl logs $pod_name) >> /host/configurator.log

    ./kubectl delete pod $pod_name
fi

echo "done [pod]; see you soon!"
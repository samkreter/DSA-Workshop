# DataScienceMicroServices
Using a simple plug and play micro service architecture with data science applications

#Deploying an AKS Cluster

1. Open up Cloudshell

    a. Go to shell.azure.com

    b. click the shell button from the Azure Portal  ![cloudshell](/images/cloudshell2.png)

2. Make sure you have the preview mode enabled
    ```Bash
    az provider register -n Microsoft.Network
    az provider register -n Microsoft.Storage
    az provider register -n Microsoft.Compute
    az provider register -n Microsoft.ContainerService
    ```

3. Create a Resource Group, see docs/AzureTerms.md for more information on what a resource group is.

    ```Bash
    az group create --name <resourceGroupName> --location eastus
    ```

4. Create the Kubernetes Cluster

    ```Bash
    az aks create --resource-group <resourceGroupName> --name <clusterName> --node-count 1 --generate-ssh-keys
    ```

    --resource-group: Name of resource group

    --name: Name for the cluster

    --node-count: How many nodes to create

    --generate-ssh-keys: We want it to handle the ssh stuff

    *NOTE:* This should show a `Running..` for a long time. That is expected. If you see an error retry it again.

5. Now we need to get access to the cluster using kubectl. To have azure set it up automagicly, run:

    ```Bash
    az aks get-credentials --resource-group <resourceGroupName> --name <clusterName>
    ```

6. Make sure everything went good by checking the nodes
    ```Bash
    kubectl get nodes
    ```

    The output should look similar to 

    ```Bash
    NAME                          STATUS    ROLES     AGE       VERSION
    k8s-myAKSCluster-36346190-0   Ready     agent     2m        v1.7.7
    ```
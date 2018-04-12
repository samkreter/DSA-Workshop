# Data Science Powered by Kubernetes and Microservices in Azure
An example of using microservice architecture powered by Kubernetes and Azure.

**NOTE:** This is a live repo. If you see something wrong or have a data science / microservice idea make a PR. I have the Contributions folder set up to add in any idea you guys have to build. I'll be checking the PRs and giving feedback. Great way to give to an open source project and get others feedback. 

# Overview

This project covers 

1. Extract, Transform and Load (ETLing). A dockerized python module will be used to
    
    a. **Extracting** data from an API, in this case bitcoin prices. Because there are a good amount of points (over a million), we will have to extract it in pieces and buffer them into the database.

    b. **Transforming** the data from the API into a readable format.

    c. **Loading** the data into a timeseries databases, [InfluxDB](https://www.influxdata.com/).

2. Data Presenation: Taking advatage of a popular dashboard called [Grafana](https://grafana.com/)

![Grafana Picture](https://github.com/samkreter/DSA-Workshop/blob/master/images/grafana.png)

# Folder Structure

* **ETL-python:** holds the code for the etling python code to get the data from a *super advanced* API and put into a database. This code even comes with a very lovely Dockerfile to dockerize the whole thing.

* **deploy:** holds the needed deployment configurations. There are configurations for both docker-compose and Kubernetes. We are much more focused on the Kubernetes here.

* **realProcessing:** code for how I got the original data. If you are interested in some golang time transforming and basic data extration using golang and python. Nothing too interesting. 

* **Contributions:** the cool stuff that you guys have come up with and decided to make a PR. Send in your ideas to show them off and get feedback!

* **api:** code for the api server if you are curious. Nothing to do with this current project. 

* **images:** just images I need to put in these markdown pages. Its way easier just to host them in the project. 

# Getting Starting

## Prerequisites

* git installed

* Docker installed

* A valid Azure Subscription

## Deploying an AKS Cluster

1. Open up Cloudshell, choose either a or b. I like a the most.

    a. Go to shell.azure.com

    b. click the shell button from the Azure Portal  
    ![](https://github.com/samkreter/DSA-Workshop/blob/master/images/cloudshell2.png)

2. Make sure you have the preview mode enabled
    ```Bash
    az provider register -n Microsoft.Network
    az provider register -n Microsoft.Storage
    az provider register -n Microsoft.Compute
    az provider register -n Microsoft.ContainerService
    ```

3. Create a Resource Group, see the [Azure Wiki Page](https://github.com/samkreter/DSA-Workshop/wiki/Azure) in this repo for more information on what a resource group is.

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

    **NOTE:** This should show a `Running..` for a long time. That is expected. If you see an error retry it again.

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

## Build the ETLing container

1. Clone the repo to your computer

    ```Bash
    git clone https://github.com/samkreter/DSA-Workshop.git
    ```

2. cd into the ETL-python folder

    ```Bash
    cd DSA-Workshop/ETL-python/
    ```

3. Build the docker image

    ```Bash
    docker build -t <dockerUsername>/dsa-etling:latest .
    ```

    **Note:** Don't miss the "." at the end. This tells docker where to look for the Dockerfile. In this case the local directory

4. Push the image to dockerhub

    ```Bash
    docker push <dockerUsername>/dsa-etling:latest
    ```

## Deploy Everything to Kubernetes Cluster

1. Back in Cloud Shell, either clone the repo again or upload the deploy/Kubernetes folder. Either drag and drop the folder onto the screen or click the upload button. 

2. Once inside the Kubernetes folder in Cloud Shell, first we are going to deploy the grafana dashboard bacuse sometimes it take Azure some time to give us a public IP address.

    ```Bash
    kubectl create -f grafana-deployment.yaml
    ```
3. Now lets deploy the Influxdb database

    ```Bash
    kubectl create -f influxdb-deployment.yaml
    ```

4. We're going nice and easy. Now lets deploy the ETL to start getting out data. Since this is just one simple container, we are going to use kubectl's run command to launch it. Make sure to replace `<dockerUsername>` with your username

    ```Bash
    kubectl run etling --image <dockerUsername>/dsa-etling:latest --restart=Never --image-pull-policy=Always
    ```

    --restart=Never=Never: if the image fails or completes, we do not want it to restart. 

    --image-pull-policy=Always: we want it to pull a new copy of the docker image down each time. Just in case we make changes. 

5. Now lets check how its doing by getting all the pods and check when its in Running state. 

    ```Bash
    kubectl get pods
    ```

6. Onces its in Running state, we can check what the logs show.

    ```Bash
    kubectl logs etling
    ```

## Setting up the Dashboard

1. By now the public ip address should have been assigned. Check the address with 

    ```Bash
    kubectl get services
    ```

Get the value for the EXTERNAL-IP value for the grafana service and pop it into a web browser.

![grafana-login](https://github.com/samkreter/DSA-Workshop/blob/master/images/grafana-login.png)

2. If you do not see a screen that looks like the one above. Something went wrong. Its over. Sorry. But if you see it you're all good. The username is **admin** and password **admin**. This is the default. 

3. Now you are taken to the setup page. Click on the create your first datasource button. 

4. Now fill in the details below

| Field | Value |
| ----- | ----- |
| Name  | influxdb |
| type  | influxdb |
| url   |  http://influxdb:8086 |
| database | BitcoinTest |
| User  | root |
| Password | root |

5. Click Save & Test. At the bottom of the page you should see a green notification saying "
Data source is working". 


## Building the Dashboard

1. Click on the **+** mark on the left and hit Dashboard. Click the Graph to start. 

2. Now an empty graph show pop up. Click the dropdown to the right of Panel Title and click edit. 

3. First on the left select influxdb from the data source dropdown. Next click the hamburger button (or stack button, just that weird looking button with 3 horizontal bars on it) and click "Toggle Edit Mode".

This changes everything to using a SQL syntax. This makes things way easier for these basic queries. 

4. Now replace the current query with 

    ```sql
    SELECT price FROM Bitcoin
    ```
5. Now when you click away, the middle of the chart should go from "No data points" to "Data points outside time range". This is just because the default is set for the last 6 hours. This is not cool enough to have data that close. 

6. Change the date range by clicking the button saying "Last 6 hours" and selecting the "Last 5 years" button. 

7. You should now be able to see the really cool curve bitcoin makes in 2017. Have fun team.


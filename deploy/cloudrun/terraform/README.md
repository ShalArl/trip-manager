# Deployment with Terraform

This directory contains Terraform configuration files to deploy the application to Google Cloud Run. Follow the instructions below to set up and deploy the application.

## Prerequisites
- [Terraform](https://www.terraform.io/downloads.html) installed on your local machine.
- A Google Cloud project with billing enabled.
- The Google Cloud SDK installed and configured on your local machine. 


## Steps to Deploy
1. **Initialize Terraform**: Navigate to the directory containing the Terraform configuration files and run the following command to initialize Terraform:
    Required to download the necessary provider plugins and set up the backend for storing the state file.

   ```bash
   terraform init
   ```
   
2. **Validate Configuration**: Run the following command to validate the Terraform configuration files for any syntax errors or issues:
   ```bash
   terraform validate
   ```

3. **Plan Deployment**: Use the following command to create an execution plan. This will show you what actions Terraform will take to deploy the application:
   ```bash
   terraform plan
    ```

4. **Apply Deployment**: If the plan looks good, run the following command to apply the changes and deploy the application to Cloud Run:
    ```bash
    terraform apply
    ```
   You will be prompted to confirm the deployment. Type `yes` to proceed.

5. **Verify Deployment**: After the deployment is complete, you can verify that the application is running on Cloud Run by visiting the URL provided in the Terraform output.


## CleanupTo clean up the resources created by Terraform, you can run the following command:
```bash
terraform destroy
```
This will destroy all the resources that were created during the deployment. You will be prompted to confirm
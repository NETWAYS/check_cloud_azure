# check_cloud_azure

check_cloud_azure is an Icinga check plugin, which is capable to check in the Microsoft Azure context.

In the current version check_cloud_azure supports the Virtual Machines context, which is capable to check a single
or multiple VMs in a defined resource group.

## Usage

### Computing - Virtual Machines

When one of the states is non-ok, or a machine is deallocated, the check will alert.

#### compute vms

Checks multiple Azure VMs in a resource group

```
Usage:
  check_cloud_azure compute vms [flags]

Flags:
  -r, --group string      Azure resource group
  -h, --help              help for vms
  -n, --tagname string    Filter resource group by tag (e.g. tag1)
  -v, --tagvalue string   Tag value of resource group (e.g. value1)

Global Flags:
      --auth-file string   Azure auth file (env:AZURE_AUTH_LOCATION)
  -s, --sub string         Azure Subscription ID (env:AZURE_SUBSCRIPTION_ID)
  -t, --timeout int        Timeout for the check (default 30)
```

```
$ check_azure_cloud compute vms --sub SUBSCRIPTION-UUID
CRITICAL - 2 VMs found - 2 running - 1 deallocated

## Group: Dev

[OK] "vm1" provision=succeeded power=running agent=succeeded
[CRITICAL] "vm2" provision=succeeded power=deallocated agent=(none)

## Group: AnotherGroup

[OK] "prod1" provision=succeeded power=running agent=succeeded
```

#### compute vm

Checks a single Azure VM

```
Usage:
  check_cloud_azure compute vm [flags]

Flags:
  -r, --group string   Azure resource group
  -h, --help           help for vm
  -n, --name string    Look for vm by name

Global Flags:
      --auth-file string   Azure auth file (env:AZURE_AUTH_LOCATION)
  -s, --sub string         Azure Subscription ID (env:AZURE_SUBSCRIPTION_ID)
  -t, --timeout int        Timeout for the check (default 30)
```

```
$ check_azure_cloud compute vms --sub SUBSCRIPTION-UUID --group group-name --name vm-name
CRITICAL - "vm-name" provision=succeeded power=deallocated agent=(none)

Size: Standard_B1s
Location: germanywestcentral
```

## Setting up Access

In order to work correctly you need the correct permissions and configuration within Azure, to grant the plugin proper
read-only access to the resources.

The following step-by-step instructions will help you to setup this configuration.

### Environment variables

The check itself needs environment variables, and supports the default environment that is compatible with other
Azure client integrations.

Export the following environment variables:

* `AZURE_TENANT_ID` See Directory Tenant ID under your APP
* `AZURE_CLIENT_ID` See Application Client ID
* `AZURE_CLIENT_SECRET` Only visible after creating a client secret for your app

Alternatively you can [create a credential file using the Azure CLI](https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli),
or manually with the following contents:

```json
{
  "tenantId": "xxx",
  "clientId": "xxx",
  "clientSecret": "xxx",
  "resourceManagerEndpointUrl": "https://management.azure.com/"
}
```

Then either set environment `AZURE_AUTH_LOCATION` or pass `--auth-file` with the file path.

Also see [Authentication methods in the Azure SDK for Go](https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization).

### App Registration

In Azure, withing the Azure Active Directory, search for the key word **App registrations** and add a new registration
with a meaningful name for the app registration like `check_cloud_aws`.

If the app registration was successfully, it should appear under the tab **Owned applications**. pen the app details
and navigate to the section **Certificates & secrets**, add a new client secret.

### Give app read access

Now the `check_cloud_azure` **App Registration** needs *read only* access to Azure to fetch monitoring values.

In Azure, search for the key word `Subscriptions`. Then click on your desired **Subscription name** and navigate to
the menu point **Access control (IAM)** and click on the button **Add role assignment**.

Select as the following:
* Role: Reader
* Assign access to: User, group, or service principal
* Select: Your_chosen_app_name

## License

Copyright (C) 2021 [NETWAYS GmbH](mailto:info@netways.)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

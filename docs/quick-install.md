<!--- OVERVIEW --->
# Quick Install

## Smart Agent Overview

The SignalFx Smart Agent is a metric-based agent written in Go that is used to monitor infrastructure and application services from a variety of environments.

The Smart Agent contains three main components:

| Component | Description |
|-----------|-------------|
| Monitors  |  This component collects metrics from the host and applications. For a list of supported monitors and their configurations, see [Monitor Configuration](./monitor-config.md).            |
| Observers |   This component collects metrics from services that are running in your environment. For a list of supported observers and their configurations, see [Observer Configuration](./observer-config.md).           |
| Writer    |   This component collects metrics from configured monitors and then sends these metrics to SignalFx on a regular basis. If you are expecting your monitors to send large volumes of metrics through a single agent, then you must update the configurations. To learn more, see [Agent Configurations](./config-schema.md#writer).          |


## Review pre-installation requirements for the Smart Agent

Before you download and install the Smart Agent on a **single** host, review the requirements below.

(For other installation options, including bulk deployments, see [Advanced Installation Options](./advanced-install-options.md).)

Please note that the Smart Agent does not support Mac OS.  

**General requirements**
- You must have access to your command line interface.
- You must uninstall or disable any previously installed collector agent from your host, such as collectd.

**Linux requirements**
- You must run kernel version 2.6 or higher for your Linux distribution.

**Windows requirements**
- You must run .Net Framework 3.5 on Windows 8 or higher.
- You must run Visual C++ Compiler for Python 2.7.

***

## Install the Smart Agent

### Step 1. Install the SignalFx Smart Agent on your host

#### Linux
<details>
<summary>Show Linux instructions</summary>
<p>
    
##### Option 1: From the SignalFx UI    

If you are reading this content from the SignalFx Smart Agent tile in the Integrations page, then simply copy and paste the following code into your command line. (The code within the tile is already populated with your realm and your organization's access token.)
    
```sh curl -sSL https://dl.signalfx.com/signalfx-agent.sh > /tmp/signalfx-agent.sh```
```sudo sh /tmp/signalfx-agent.sh --realm YOUR_SIGNALFX_REALM YOUR_SIGNALFX_API_TOKEN```

***

##### Option 2: From the documentation site 

If you are reading this content from the SignalFx documentation site, then SignalFx recommends that you access the Integrations page in the SignalFx UI to copy the pre-populated installation code.  

1. Log in to SignalFx and click the **Integrations** tab to open the Integrations page. Look for the SignalFx Smart Agent tile. You can search for it by name, or find it in the **Essential Services** section.
2. Under **Essential Services**, click **SignalFx Smart Agent**.
3. Click **Setup**.
4. Locate the text box for Linux users.
5. Copy, paste, and run the code in your command line. (The code within the tile is already populated with your realm and your organization's access token.)  

</p>
</details>

***

#### Windows


<summary>Show Windows instructions</summary>
<p>

##### Option 1: From the SignalFx UI    
If you are reading this content from the SignalFx Smart Agent tile in the Integrations page, then simply copy and paste the following code into your command line. (The code within the tile is already populated with your realm and your organization's access token.)


```sh
& {Set-ExecutionPolicy Bypass -Scope Process -Force; $script = ((New-Object System.Net.WebClient).DownloadString('https://dl.signalfx.com/signalfx-agent.ps1')); $params = @{access_token = "YOUR_SIGNALFX_API_TOKEN"; ingest_url = "https://ingest.YOUR_SIGNALFX_REALM.signalfx.com"; api_url = "https://api.YOUR_SIGNALFX_REALM.signalfx.com"}; Invoke-Command -ScriptBlock ([scriptblock]::Create(". {$script} $(&{$args} @params)"))}
```

***

##### Option 2: From the documentation site 
If you are reading this content from the SignalFx documentation site, then SignalFx recommends that you access the Integrations page in the SignalFx UI to copy the pre-populated installation code.  

1. Log in to SignalFx and click the **Integrations** tab to open the Integrations page. Look for the SignalFx Smart Agent tile. You can search for it by name, or find it in the **Essential Services** section.
2. Under **Essential Services**, click **SignalFx Smart Agent**.
3. Click **Setup**.
4. Locate the text box for Windows users.
5. Copy, paste, and run the code in your command line. (The code within the tile is already populated with your realm and your organization's access token.)  


The agent will be installed as a Windows service and will log to the Windows Event Log.
</p>
</details>

***


### Step 2. Confirm your Installation


1. To confirm your installation, enter the following command on the Linux or Windows command line: 

    ```sh
    sudo signalfx-agent status
    ```

    The return should be similar to the following example:  

    ```sh
    SignalFx Agent version:           4.7.6
    Agent uptime:                     8m44s
    Observers active:                 host
    Active Monitors:                  16
    Configured Monitors:              33
    Discovered Endpoint Count:        6
    Bad Monitor Config:               None
    Global Dimensions:                {host: my-host-1}
    Datapoints sent (last minute):    1614
    Events Sent (last minute):        0
    Trace Spans Sent (last minute):   0
    ```

2. To confirm your installation, enter the following command on the Linux or Windows command line: 

    | Command | Description   |
    |---|---|
    | <code>signalfx-agent status config</code>   | This command shows resolved config in use by the Smart Agent. |
    | <code>signalfx-agent status endpoints</code>  | This command shows discovered endpoints.  |
    | <code>signalfx-agent status monitors</code>  | This command shows active monitors.  |
    | <code>signalfx-agent status all</code>  | This command shows all of the above statuses. |

***

### Troubleshoot the Smart Agent installation

If you are unable to install the Smart Agent, consider reviewing your error logs: 

For Linux, use the following command to view error logs via Journal:

```sh
journalctl -u signalfx-agent | tail -100
```

For Windows, review the event logs.

***

For additional installation troubleshooting information, including how to review logs, see [Frequently Asked Questions](./faq.md).

***

### Review additional documentation

After a successful installation, learn more about the SignalFx agent and the SignalFx UI. 

* Review the capabilities of the SignalFx Smart Agent. See [Advanced Installation Options](./advanced-install-options.md).

* Learn how data is displayed in the SignalFx UI. See [View infrastructure status](https://docs.signalfx.com/en/latest/getting-started/quick-start.html#step-3-view-infrastructure-status).

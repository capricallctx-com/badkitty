# Bad Kitty
## Description

`Bad Kitty` is a replacement web server for container stacks and websites, that
aims to not be the fastest, nor the most featureful, but the most secure and best suited for modern workloads.  `Bad Kitty` doesn't do
everything, but it does most of the things regular people need simply and securely, and to support _reactive_ configuration changes.



No plugins, no modules, no extra configuration.  Just a simple, secure web server.

## WTF is 'reactive' configuration changes?

`Bad Kitty` is designed to be configured by a simple HCL file for basic functions.  Part of this configuration is the ability to connect to a _YOWL server_ either on the cloud or on-premises.

_YOWL_, if enabled, pushes insights from your kitty colony to a central location, and allows hints to be pushed into the colony members.  For the open source _YOWL_ - humans curate the hints to prevent shenanigans.   For private stacks - you can probably just let them ride.

### What does this mean?

_YOWL_ data allows the colony to react to ongoing threats and changes.   Bad actors and port scanners can be shared among the colony.
When an IP address triggers a block for scanning for known vulnerabilities, the data goes into _YOWL_ and it can be automatically or
manually put into timeout.  _YOWL_ is aware of any security flaws in `Bad Kitty` releases and can warn you or report coverage.

If not manually configured, this also restricts broken or insecure protocols and keeps your pack consistent and secure.
If a new issue with a protocol (like ChaCha20) crops up - the colony can disable it immediately.

You can run your own _YOWL_, or use the community one.  Community _YOWL_ rule changes are run through our Github and are transparent to the users.   And you can always just disable _YOWL_ if you don't want to participate.



This allows `Bad Kitty` to receive configuration changes in real time, and apply them without restarting the server.  This allows for a more secure and responsive server, and allows for more complex configurations to be managed by a central server.

## Features

| Feature                           | Description                                                                                                                                                                                                                                                                                     | Status  |
|-----------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------|
| Proactive Vulnerability Reporting | _YOWL_ will report issues with your install to your logs or a local SMTP server.  Don't wait for scans or reports when it's too late.                                                                                                                                                           | Present |
| HTTP or HTTPS Hosting             | `Bad Kitty` can host your website or web app over HTTP or HTTPS.  Default is redirecting HTTP to HTTPS (the right thing to do).                                                                                                                                                                 | Present |
| ACME Out of the Box | Huzzah!  No more downloading shady libraries or fiddling.  `Bad Kitty` just works.  Why does everything not do this?                                                                                                                                                                            | Planned |
| Best Practices Enforcement        | `Bad Kitty` can enforce best practices for SSL, including HSTS, and other security headers automatically.  Don't spend hours fiddling with disabling or enabling protocols - `Bad Kitty` does it automatically. It also sets the CORS and other headers to protect against most common attacks. | Planned |
| Managed Static Content            | `Bad Kitty` can serve your web app for most cases.  Built your Angular, React, or plain JS app and just add the `static`yaml stanza to serve it                                                                                                                                                 | Present |
| Fuzz Buster                       | `Bad Kitty` can automatically block IP addresses that are scanning your site for vulnerabilities.  It can also block IP addresses that are scanning for vulnerabilities on other sites if _YOWL_ is used.                                                                                       | Planned |


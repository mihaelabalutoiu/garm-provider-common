# GARM provider common

This package contains common code for GARM providers and for GARM itself.

## Generating cloud userdata

One of the main functions of this package is to allow you to generate userdata scripts for various cloud providers. In most (if not all) cases, you will need a runner installation script and some cloud specific way to wrap it in. Most clouds will support `cloud-init` userdata for most Linux distributions. There are a few outliers that use a different initialization system like [ignition](https://github.com/coreos/ignition), but for those cases, a potential provider will be able to easily wrap the install script in a config suitable for whatever init system is used.

### Userdata

The [cloudconfig](./cloudconfig) package in this repo gives you a set of helper functions that allow you to generate a suitable install script that should work with GARM. I will give a few details in regards to how to use it.

Generating the userdata for your target OS/Cloud combination implies a few steps:

* Generate the install script
* Add any SSH keys to the config
* Add any extra scripts
* Generate the final userdata

For most casea, you can simply run [GetCloudConfig()](https://github.com/cloudbase/garm-provider-common/blob/main/cloudconfig/util.go#L176) which will generate a `cloud-init` cloud config for Linux and a powershell script for Windows, using the default templates and template context we supply.

You may, however override the default install scripts if you wish. The [same context](https://github.com/cloudbase/garm-provider-common/blob/main/cloudconfig/templates.go#L418-L458) we use to generate the install script using the default templates, will be passes into your template. However, your script may need aditional information that we don't supply. For situations like this, we have an [ExtraSpecs](https://github.com/cloudbase/garm/blob/main/doc/extra_specs.md) field that is specific to userdata to help out here. This is what the [ExtraContext](https://github.com/cloudbase/garm-provider-common/blob/main/cloudconfig/templates.go#L454-L457) field is for.

Overriding the default runner installation script is as easy as creating an extra specs json with the following contents:

```json
{
    "runner_install_template": "BASE64_ENCODED_TEMPLATE",
    "extra_context": {
        "extra_key": "extra_value"
    },
}
```

and updating your pool to use it:

```bash
garm-cli pool update --extra-specs-file extra_specs.json <POOL ID>
```

Note: If you override the default template, it falls onto you to ensure the correctness and suitability of this template for your target OS/Cloud combination.

With these options set, calling `GetCloudConfig()` will use your template instead of the default one. You still get a `cloud-init` config for Linux using this function. So what do we do if we need more granular control over how userdata is generated?

The [cloudconfig](./cloudconfig) package exposes a few more functions that allow you to generate the install script and the cloud config separately. The biggest chunk of the userdata script is the actual install script which is added as a file and then executed by `cloud-init`. But as we mentioned, you may use a different cloud initialization system. To generate just the install script, you can call the [GetRunnerInstallScript()](https://github.com/cloudbase/garm-provider-common/blob/main/cloudconfig/util.go#L74) function, directly. Have a look at the package for more details.
# DevStream Security Policy

Version: **v0.1 (2022-02-28)**

## Overview

As a tool manager, some of the tools need to be installed in a production environment. Thus, some of DevStream's plugins need to have production access, which makes security an important topic.

The DevStream team takes security as job zero and improves it continuously. 

## Supported Versions

We currently support the most recent release (`N`, e.g. `0.2`) and the release previous to the most recent one (`N-1`, e.g. `0.1`). With the release of `N+1`, `N-1` drops out of support and `N` becomes `N-1`.

We regularly perform patch releases (e.g. `0.1.1` and `0.2.1`) for supported versions, which will contain fixes for security vulnerabilities. Prior releases might receive critical security fixes on a best effort basis, however, it cannot be guaranteed that security fixes get back-ported to these unsupported versions.

In some cases, where a security fix needs complex re-design of a feature or is otherwise very intrusive, and there's a workaround available, we may decide to provide a forward-fix only, e.g. to be released the next minor release, instead of releasing it within a patch branch for the currently supported releases.

## Reporting a Vulnerability

Please report vulnerabilities by e-mail to the following address: 

- tiexin.guo@merico.dev
- tao.hu@merico.dev
- fangbao.li@merico.dev

If you find a security-related bug in DevStream, we kindly ask you to disclose responsibly and give us appropriate time to react to mitigate the vulnerability.

We do our best to react quickly. Sometimes, it might take a little longer (e.g. out-of-office conditions); please bear with us in these cases.

We publish security advisories using the [Git Hub Security Advisories](https://github.com/devstream-io/devstream/security/advisories) feature to keep our community well informed, and will credit you for your
findings (unless you prefer to stay anonymous, of course).

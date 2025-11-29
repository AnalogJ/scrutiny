# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| Latest  | Yes                |
| < Latest | No                |

Only the latest release receives security updates. We recommend always running the most recent version.

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report security vulnerabilities by emailing the maintainer directly. You can find contact information on the maintainer's GitHub profile: [@Starosdev](https://github.com/Starosdev)

Please include the following information in your report:

- Type of vulnerability (e.g., SQL injection, XSS, authentication bypass)
- Full paths of source file(s) related to the vulnerability
- Location of the affected source code (tag/branch/commit or direct URL)
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it

## Response Timeline

- **Initial Response**: Within 72 hours of report submission
- **Status Update**: Within 7 days with an assessment of the vulnerability
- **Resolution**: Security patches are prioritized and typically released within 30 days for confirmed vulnerabilities

## Disclosure Policy

- We will acknowledge receipt of your vulnerability report
- We will confirm the vulnerability and determine its impact
- We will release a fix as soon as possible, depending on complexity
- We will publicly disclose the vulnerability after a fix is available

## Security Best Practices for Users

1. **Keep Updated**: Always run the latest version of Scrutiny
2. **Network Security**: Do not expose the Scrutiny web interface directly to the internet without proper authentication
3. **Docker Security**: Follow Docker security best practices when deploying containers
4. **File Permissions**: Ensure configuration files containing sensitive data have appropriate permissions

## Scope

This security policy applies to:

- The Scrutiny web application
- The Scrutiny collector
- Official Docker images from this repository

Third-party integrations, forks, and unofficial distributions are outside the scope of this policy.

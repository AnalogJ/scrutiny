# Ansible Install

[Zorlin](https://github.com/Zorlin) has developed and now maintains [an Ansible playbook](https://github.com/Zorlin/scrutiny-playbook) which automates the steps involved in manually setting up Scrutiny.

Using it is simple:

* Grab a copy of the playbook
* Follow the directions in the playbook repository
* Run `ansible-playbook site.yml`
* Visit http://your-machine:8080 to see your new Scrutiny installation.

It will automatically pull metrics from machines once a day, at 1am.

You can see it in action below.

[![asciicast](https://asciinema.org/a/493531.svg)](https://asciinema.org/a/493531)


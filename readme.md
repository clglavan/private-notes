# Private Notes - send self-distructing notes over the internet

![private_notes](private_notes.gif)

Send private notes over the internet as one time links that destroy themselves after they are read. Optionally chose a password for your note.

This repo wishes to provide an open-source alternative for managed solutions of similar usecase. Advantages include future possibility for branding and contents of messages being saved on the organization's own cloud resources.

# Getting started


## Deployments

Deployment config, coming soon.

# Run locally

Check contents of `local.sh` 

# Known bugs
- [ ] if you open a password-protected note, trying to decrypt with a wrong password will not work ( as expected ) but will also trigger the note to be destroyed.
# Further improvements
- [x] make decryption by choice, with "view note" button
- [x] add custom password
- [x] refactor html with layout
- [ ] terraform code for deployment cloud run
- [ ] fix cf deployment after merge
- [ ] make config.yaml for easier app behaviour tweaks and branding
- [ ] enable easy custom branding
- [ ] add copy to clipboard button for secret link
- [ ] refactor code
    - [ ] implement routing
    - [ ] implement middlewares
    - [ ] implement logging

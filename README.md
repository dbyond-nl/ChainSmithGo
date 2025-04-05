
# ChainSmith

A tool to quickly, easily and reproducibly generate Cryptographic certificates for use
with PostgreSQL cluster servers.

## Usage

Chainsmith currently uses a 'configuration file' to provide it with the required information. For more information on the format of that file see (TODO: insert link)

chainsmith -c /PATH/TO/CONFIG/chainsmith.yml

For more options, see

chainsmith --help

## Installation

Please see [QUICKSTART.md] for options to install and run chainsmith.

## Configuration file

An example config file chainsmith.yml is shipped with chainsmith. Change as required and run chainsmith with the '-c' option and a Yaml (.yml) file.

## Output format

Note that by default the certificates are written as a yaml hash to stdout, and the private keys are written as a yal hash to stderr. Alternatively you can redirect them to files using the -o and -p options.

# Background

## Origin

Sebastiaan Mannem developed Chainsmith (originally in Python) as part of a set of tools he created to make PostgreSQL easier to use. It was converted to Golang and was made part of the pgVillage tool set to give it wider reach.

## Why use ChainSmith

To run Enterprise services securely you should use SSL for encryption in transit and for verifying trust between systems. Unfortunately, creating a simple chain with a root certificate, 2 intermediate certificates and client certificates and/or server certificates using standard tooling is a very complex procedure requirying much manual effort and knowledge. This project is meant to reduce the needed effort and knowledge without compromising the security of the process.

With ChainSmith, you can define a chain in yaml config, and then run this script to create a root CA, intermediates and signed certificates.

All tar files are bundled in separate yaml files, so you can easily use them in tools like Ansible for deployment. Or, if you do want externally signed certificates, you can use ChainSmith to generate all Certificate Signing Requests (CSR's) to be signed externally. And you can run with the generated chain until the externally signed certificates are available.

ChainSmith is a crucial piece into improving adoption of running Postgres and other tools with proper security. And as such systems can be easily equipped with the proper certificate chains so that secure communication and authorization is possible.

## Why you should use certificates

Certificates are a technical implementation for verification of trustworthiness. Certificates can be verified on the following points:

    to be used for its correct purpose
    to be used by the correct person or system
    to be used by a person or system which is trusted by you, or a party you trust Once trustworthiness is established, certificates can be used to limit communication to only the 2 parties that are communicating.

### How verification of trust works

A certificate can be verified to:

- be used for its proper purpose
  - common name should correspond to the server you are communicating with, or
  - common name should correspond to the user trying to authenticate with it
- be used by the proper system or user
  - the certificate can be shared to everyone that wants to be verified, but
  - the certificate can only be used by those that hold the corresponding private key
- be handed out by someone or something you trust, or someone they trust
  - Every certificate is signed by another certificate (except for root certificates)
  - Before signing off on a certificate, the authority is required to properly verify that the certificate is requested by the proper person, system or authority. This creates a chain of trust and if you can trust one certificate in the chain, you can also trust all that are signed by that certificate

Other properties of Certficates are:

- Certificates that can no longer be trusted can be revoked
- Once trust is verified, communication is assured to be protected from anyone besides the two parties that are communicating
- All information encrypted with the certificate can only be decrypted by the system or person with the corresponding private key

# Development

This project is maintained on github.

If you run into issues while using, or you may have other suggestions to improve ChainSmith, please create an Issue.

And if you want to contribute, don't be shy, just create a Pull Request and we will probably merge.
License
This software (all code in this github project) is subjective to GNU GENERAL PUBLIC LICENSE version 3.

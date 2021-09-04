# README #

DUM - is a pet-project system for monitoring missing updates on the windows machines.

### What is this repository for? ###

DUM relies on different set of data:
* Windows Update service data of specific machine
* WSUS deployment data
* May be any other in the future

So as this project is a pet-project, the author will try to use such patterns and principles like:
* Clean architecture
* Domain-Driven Development
* Microservices
* CQRS
* SOLID
* GRASP

The practical scenarios that system should support are displayed at the following diagram:

![Cases](/diagrams/cases.jpg?raw=true)

The system consists of the following components:

![Components](/diagrams/components.jpg?raw=true)

Components:
* MachineService - service, allows to perform an operations with specific machine. Provides the ReportAPI interface, which is used for reporting about any updates, that are missing on the specific machine.

Go modules:

![Modules](/diagrams/modules.jpg?raw=true)

Modules:
* Machines - module, that implements MachineAPI component itself
* Contracts/Machines - module with set of dto, which are needed for using MachinesService

### How do I get set up? ###

All you need is to follow the instructions for each microservice. It can be found at README.md of every directory with microservice code

### Contribution guidelines ###

There is no contribution guideline yet.
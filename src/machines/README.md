# Machines Microservice #

MachinesService - is a microservice, responbible only for machines domain entity.
Basically, it allows to tell the DUM that some machine has some missing updates.
This request will be processed in such a way:

![Machine report sequence](/diagrams/machine_report.jpg)

From the sequence diagram, we can form the 4-layered microservice, according to Clean architecture:

![Machines classes](/diagrams/machine_module.jpg)

To create Docker image, please use make_image.bat.
To run microservice just build it and run the executable - no specific setup is needed.
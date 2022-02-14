# Deploys

*NOTE* This code was built assuming terraform version >= 1.

Demonstrating an automated deploy pipeline for a toy project is nuanced. 

For something as simple as lakitu, AWS Lightsail would likely suffice (the modern successor to Elastic Beanstalk)
as it now support containers. It's features are quite limited and it abstracts away almost everything
in order to provide the lowest barrier to entry. It's certainly not what I would recommend for production use, however.
The absolute _best_ tool would be Kubernetes, on AWS EKS, but that requires a much larger investment
in both cost and upfront planning, and would be overkill for a toy project demonstration. A good middle ground
would be AWS Fargate; it allows for multiple deploy strategies, finer-grained control of container resources and
flexibility with regards to other service integrations.

Even with the above in mind, creating the necessary AWS resources to support a solution would be difficult to
demonstrate "live"; it would require an AWS development account, real dollars to launch and test, and additional
security integrations and controls to provide an end-to-end solution, which is out of scope of this exercise.

The terraform resources here create an OIDC provider and IAM Role to allow for Github to push built images.
This happens automatically on a push to `main` with a new `VERSION`.

The application of these resources was done locally; in a production setting this would also be automated, 
either through Terraform Cloud or a kubernetes operator. It uses a local state store, which would need to changed
to a remote backend, with resource access policies implemented to make sure infrastructure changes are secure.

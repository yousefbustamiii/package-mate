package components

var catalogContainers = Section{
	Name: "Containers & DevOps",
	Items: []InstallItem{
		{
			Name:   "Docker",
			Desc:   "Platform for building and running containerized applications",
			Cask:   "docker",
			Color:  rgb(80, 180, 255),
			Binary: "docker",
		},
		{
			Name:    "Docker Compose",
			Desc:    "Tool for defining and running multi-container Docker systems",
			Formula: "docker-compose",
			Color:   rgb(70, 140, 255),
			Binary:  "docker-compose",
		},
		{
			Name:    "Colima",
			Desc:    "Container runtime for macOS with minimal resource overhead",
			Formula: "colima",
			Color:   rgb(50, 190, 230),
			Binary:  "colima",
		},
		{
			Name:    "Kubernetes (kubectl)",
			Desc:    "Command-line tool for controlling Kubernetes system clusters",
			Formula: "kubernetes-cli",
			Color:   rgb(100, 160, 255),
			Binary:  "kubectl",
		},
		{
			Name:    "k9s",
			Desc:    "Terminal-based UI for managing Kubernetes cluster operations",
			Formula: "k9s",
			Color:   rgb(255, 210, 50),
			Binary:  "k9s",
		},
		{
			Name:    "Helm",
			Desc:    "Package manager for managing Kubernetes application charts",
			Formula: "helm",
			Color:   rgb(50, 130, 200),
			Binary:  "helm",
		},
		{
			Name:    "Terraform",
			Desc:    "Infrastructure as Code tool for managing cloud resources",
			Formula: "terraform",
			Color:   rgb(120, 80, 230),
			Binary:  "terraform",
		},
		{
			Name:    "Ansible",
			Desc:    "Automation platform for IT configuration and management",
			Formula: "ansible",
			Color:   rgb(255, 100, 50),
			Binary:  "ansible-playbook",
		},
		{
			Name:    "awscli",
			Desc:    "Universal interface for Amazon Web Services cloud infrastructure",
			Formula: "awscli",
			Color:   rgb(255, 170, 50),
			Binary:  "aws",
		},
		{
			Name:   "gcloud CLI",
			Desc:   "Command-line interface for Google Cloud Platform services",
			Cask:   "google-cloud-sdk",
			Color:  rgb(100, 170, 255),
			Binary: "gcloud",
		},
		{
			Name:    "Azure CLI",
			Desc:    "Command-line experience for Microsoft Azure cloud management",
			Formula: "azure-cli",
			Color:   rgb(50, 160, 230),
			Binary:  "az",
		},
	},
}

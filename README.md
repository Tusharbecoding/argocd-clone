# GitOps Sync Tool

An open-source GitOps tool inspired by ArgoCD, designed to automate the deployment of Kubernetes resources from a Git repository. This project enables automatic synchronization and application of Kubernetes manifests, making it easy to manage your cluster state as code.

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Setup Guide](#setup-guide)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [Contributing](#contributing)

---

## Features

- **Automated Sync**: Periodically syncs with a Git repository to detect and apply changes to Kubernetes resources.
- **GitOps-Friendly**: Supports automated and declarative management of Kubernetes manifests.
- **Customizable Poll Interval**: Configurable polling interval to check for Git updates as frequently as you need.

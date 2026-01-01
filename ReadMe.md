# Distributed Key-Value Store

A horizontally scalable, fault-tolerant distributed Key-Value store with sharding, replication, and centralized orchestration.

---

## 🧠 Architecture Overview

This system consists of:
- A **Master Server** responsible for routing, sharding, and cluster metadata
- Multiple **Storage Nodes** responsible for durable data storage
- Clients interact with the system via a unified API Gateway

---

## 👑 Master Node Architecture

```mermaid
graph TD
    %% External Client actor
    Client[👤 External Client]

    %% The Master Server Subsystem
    subgraph "👑 Master Server (Orchestrator)"
        direction TB
        APIGateway["🚪 API Gateway / Load Balancer"]
        ShardingLogic["🧮 Sharding & Router Logic<br/>(Hash(key) → Node IP)"]
        MetadataStore[("📒 Cluster Metadata<br/>(Node Inventory & Health)")]
    end

    %% The Cluster of Storage Nodes
    subgraph "🧱 Storage Cluster (Data Nodes)"
        direction LR
        Node1[("📦 Node 1<br/>(:3001)")]
        Node2[("📦 Node 2<br/>(:3002)")]
        Node3[("📦 Node 3<br/>(:3003)")]
    end

    %% Flows
    Client -->|"1. PUT /key (Save Data)"| APIGateway
    APIGateway --> ShardingLogic
    ShardingLogic -.->|"Lookup Active Nodes"| MetadataStore
    
    %% Replication Flow (Fan-out)
    ShardingLogic ==>"2. Stream Data (Replica 1)"==> Node1
    ShardingLogic ==>"2. Stream Data (Replica 2)"==> Node2
    ShardingLogic -.->|"2. Stream Data (Replica 3 - Optional)"| Node3

    %% Node Independence
    Node1 --- Disk1[(Disk 1)]
    Node2 --- Disk2[(Disk 2)]
    Node3 --- Disk3[(Disk 3)]

    style APIGateway fill:#e3f2fd,stroke:#1565c0,color:#000
    style ShardingLogic fill:#fff9c4,stroke:#fbc02d,color:#000
    style MetadataStore fill:#fce4ec,stroke:#c2185b,color:#000
    style Node1 fill:#e8f5e9,stroke:#2e7d32,color:#000
    style Node2 fill:#e8f5e9,stroke:#2e7d32,color:#000
    style Node3 fill:#e8f5e9,stroke:#2e7d32,color:#000
```

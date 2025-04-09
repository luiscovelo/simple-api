# Infraestrutura do cluster

# Control Plane

Conjunto de aplicações complexas responsável por gerenciar todo o ciclo de vida dos objetos dentro do Cluster.

Exige muito cuidado e conhecimento para manter e gerenciar, pois qualquer problema dentre essas aplicações todo o seu cluster pode parar de funcionar.

É composto pelas seguintes aplicações:

- etcd
    - É o banco de dados do cluster, todo o conteúdo e estado é armazenado nele, seu armazenamento é por meio de Chave=Valor.
    - Se houver um problema/corrompimento de dados seu cluster provavelmente irá parar.
- scheduler
    - Responsável por escolher em Node um POD deverá ser criado.
    - Leva em consideração os recursos disponíveis, politicas, etc, para tomar a decisão.
- controller-manager
    - Responsável por ficar olhando para o estado do cluster atual e tenta mover para o estado desejado por meio da API Server.
- api-server
    - Expõe uma API para outras partes do cluster comunicar com o control-plane, tudo passa por essa API.
- cloud-controller-manager
    - Componente que permite conectar seu cluster à api do cloud provider, sendo a ponte entre o cloud-provider -> cloud-controller-manager -> API Server.

# Control Data

## Node ou Worker Nodes

São máquinas virtuais ou físicas onde os containers serão executados.

Representa os recursos computacionais do cluster

## Kubelet

É um agente que é executado em cada Node do cluster, assegura que os containers estão rodando em um POD.

## Kube-proxy

É um proxy de rede que roda em cada Node do cluster.

Mantém regras de rede dentro do Node, isso permite a comunicação entre Pods dentro e fora do cluster.

## Container runtime

Responsável por gerenciar a execução e o ciclo de vida do container.

Kubernetes suporta: containerd, CRI-O, entre outros que implementam a interface CRI (Container Runtime Interface)

---

# Workloads

## POD

É o menor objeto/unidade implantável dentro do K8S.

Um pod pode container 1 ou mais containers mas geralmente contém apenas 1.

Só faz sentido ter mais de 1 container dentro do Pod se houver uma relação intrinsica entre eles, por exemplo um banco de dados com um agente de backup, o agente de backup depende do cliente do banco, desta forma faz sentido ter rodando dentro de um mesmo pod.

# Tipos de gerenciamento de workload

## Deployment

Deployment é um objeto que agrega um conjunto de regras e definições para criar os pods, por exemplo:

Gostaria de 10 replicas de containers rodando a imagem “my-api:latest” com 1Gi de memória e 2vCPU, expondo a porta 8080.

## ReplicaSet

ReplicaSet é um objeto que garante que as regras e definições definida no Deployment sejam aplicadas, gerenciando todo o ciclo de vida dos PODS.

Ela é quem irá garantir que existirá sempre 10 pods rodando com a imagem “my-api” na porta 8080, em caso de crash ou não resposta do pod, o ReplicaSet restarta, deleta ou cria o Pod para normalizar e obdecer as regras e definições.

Ou seja, em ordem de ligação ficaria assim: Deployment → ReplicaSet → POD

---

# Controllers

## Service

Services são abstrações que definem um conjunto lógico de Pods e uma política para acessá-los, funcionando como uma camada de abstração de rede.

É por meio deste objeto que é possível a comunicação entre pods, por exemplo:

Considere 2 pods em execução chamados `api` e `report-api`, cada um está com service configurado com os respectivos nomes: `api-service` e `report-api-service`, desta forma o pod `api` consegue facilmente realizar uma chamada http para o `report-api` da seguinte forma:

API →POST http://report-api-service:8080/generate-report

O service possui alguns tipos de comunicação, sendo:

- NodePort
    - Expõe a porta em cada node contido dentro do cluster. Ficando acessível externamente através do IP do node.
    - Range de portas permitidas: 30000-32767
- ClusterIP
    - Este tipo permite comunicação entre serviços apenas por dentro do cluster ou seja suas aplicações só podem ser acessadas via objeto externo como um Load Balancer ou Proxy
- LoadBalancer
    - Este tipo expõe o serviço externamente usando o balanceador de carga do provedor de nuvem, criando automaticamente um NodePort e um ClusterIP.
- ExternalName
    - Mapeia o serviço para um nome DNS especificado, sem usar seletores de Pod ou definir endpoints.
- Headless Service
    - Um tipo especial de ClusterIP onde não há IP de cluster (clusterIP: None), útil quando você precisa descobrir diretamente os IPs dos Pods individuais.

## Ingress

Ingress expõe rotas HTTP/HTTPS de fora do cluster para os serviços dentro do cluster.
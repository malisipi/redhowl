let template_agents = document.head.querySelector("template.agents");
let view_agents_list = document.body.querySelector(".agents-list");

let update_agents = async () => {
    let req = await fetch("/api/agents");
    let body = await req.json();
    
    let agents_list = Array.from(view_agents_list.childNodes);
    let agents_uuid_list = agents_list.map(a=>a.getAttribute("uuid"));

    for(let i=0; i<body.agents.length; i++){
        let current_agent = body.agents[i];
        let agent_uuid = current_agent.uuid;

        let agent_index = agents_uuid_list.indexOf(agent_uuid);
        let agent_node = null;
        if(agent_index == -1){
            agent_fragment = template_agents.content.cloneNode(true);
            agent_node = agent_fragment.firstElementChild;

            // static infos
            agent_node.setAttribute("uuid", agent_uuid);
            if(current_agent.metrics.os.generic == "darwin"){
                agent_node.querySelector(".logo-os").setAttribute("os", "darwin");
            } else if(current_agent.metrics.os.generic == "linux") {
                let platform = current_agent.metrics.os.platform.toLowerCase();
                if(platform.includes("cachy")){
                    agent_node.querySelector(".logo-os").setAttribute("os", "cachy");
                } else if(platform.includes("debian")){
                    agent_node.querySelector(".logo-os").setAttribute("os", "debian");
                } else if(platform.includes("ubuntu")){
                    agent_node.querySelector(".logo-os").setAttribute("os", "ubuntu");
                } else if(platform.includes("fedora")){
                    agent_node.querySelector(".logo-os").setAttribute("os", "fedora");
                } else if(platform.includes("arch")){
                    agent_node.querySelector(".logo-os").setAttribute("os", "arch");
                } else {
                    agent_node.querySelector(".logo-os").setAttribute("os", "linux");
                };
            } else if(current_agent.metrics.os.generic == "android") {
                agent_node.querySelector(".logo-os").setAttribute("os", "android");
            } else if(current_agent.metrics.os.generic == "windows") {
                if(current_agent.metrics.os.platform.toLowerCase() == "wine"){
                    agent_node.querySelector(".logo-os").setAttribute("os", "wine");
                } else {
                    agent_node.querySelector(".logo-os").setAttribute("os", "windows");
                };
            } else {
                agent_node.querySelector(".logo-os").setAttribute("os", "unknown");
            }
            agent_node.querySelector(".logo-arch").innerText = current_agent.metrics.os.arch;

            agent_node.querySelector("button.authorize").addEventListener("click", function (event, _agent_uuid=agent_uuid){
                fetch("/api/agents/authorize", {method:"POST", body:JSON.stringify({uuid:_agent_uuid})})
            });
            agent_node.querySelector("button.unauthorize").addEventListener("click", function (event, _agent_uuid=agent_uuid){
                fetch("/api/agents/unauthorize", {method:"POST", body:JSON.stringify({uuid:_agent_uuid})})
            });

            view_agents_list.append(agent_node);
        } else {
            agent_node = agents_list[agent_index];
        };

        let meter_cpu = agent_node.querySelector(".meter.cpu");
        meter_cpu.querySelector("progress").value = current_agent.metrics.cpu;
        meter_cpu.querySelector(".info").innerText = (current_agent.metrics.cpu * 100).toFixed(1) + "%";

        let meter_mem = agent_node.querySelector(".meter.memory");
        meter_mem.querySelector("progress").value = (current_agent.metrics.memory.used/(current_agent.metrics.memory.total != 0 ? current_agent.metrics.memory.total : 1));
        meter_mem.querySelector(".info").innerText = current_agent.metrics.memory.used.toFixed(1) + "/" + current_agent.metrics.memory.total.toFixed(1) + " GiB";

        let meter_disk = agent_node.querySelector(".meter.disk");
        meter_disk.querySelector("progress").value = (current_agent.metrics.disk.used/(current_agent.metrics.disk.total != 0 ? current_agent.metrics.disk.total : 1));
        meter_disk.querySelector(".mount-point").innerText = (current_agent.metrics.disk.mountPoint != "") ? "(" + current_agent.metrics.disk.mountPoint + ")" : "";
        meter_disk.querySelector(".info").innerText = current_agent.metrics.disk.used.toFixed(1) + "/" + current_agent.metrics.disk.total.toFixed(1) + " GiB";

        agent_node.querySelector(".authorization").innerText = current_agent.status;
    };
};

update_agents();

setInterval(update_agents, 1000);
//setTimeout(()=>location.reload(), 1000);
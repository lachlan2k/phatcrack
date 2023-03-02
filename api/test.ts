/* Do not change, this code is generated from Golang structs */


export class AgentCreateResponseDTO {
    name: string;
    id: string;
    key: string;

    constructor(source: any = {}) {
        if ('string' === typeof source) source = JSON.parse(source);
        this.name = source["name"];
        this.id = source["id"];
        this.key = source["key"];
    }
}
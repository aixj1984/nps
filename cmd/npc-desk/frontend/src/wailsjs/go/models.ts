export namespace configdoc {
	
	export class Section {
	    id: string;
	    title: string;
	    name: string;
	    enabled: boolean;
	    kind: string;
	    fields: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new Section(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.name = source["name"];
	        this.enabled = source["enabled"];
	        this.kind = source["kind"];
	        this.fields = source["fields"];
	    }
	}
	export class Document {
	    sections: Section[];
	    mode: string;
	    warnings?: string[];
	    errors?: string[];
	
	    static createFrom(source: any = {}) {
	        return new Document(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sections = this.convertValues(source["sections"], Section);
	        this.mode = source["mode"];
	        this.warnings = source["warnings"];
	        this.errors = source["errors"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace npc {
	
	export class CommonSettings {
	    server_addr: string;
	    conn_type: string;
	    vkey: string;
	    dns_server: string;
	    auto_reconnection: boolean;
	    local_ip: string;
	    ntp_server: string;
	    ntp_interval: number;
	    max_conn: number;
	    flow_limit: number;
	    rate_limit: number;
	    crypt: boolean;
	    compress: boolean;
	    tls_enable: boolean;
	    proxy_url: string;
	    disconnect_timeout: number;
	    pprof_addr: string;
	    web_username: string;
	    web_password: string;
	    basic_username: string;
	    basic_password: string;
	
	    static createFrom(source: any = {}) {
	        return new CommonSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.server_addr = source["server_addr"];
	        this.conn_type = source["conn_type"];
	        this.vkey = source["vkey"];
	        this.dns_server = source["dns_server"];
	        this.auto_reconnection = source["auto_reconnection"];
	        this.local_ip = source["local_ip"];
	        this.ntp_server = source["ntp_server"];
	        this.ntp_interval = source["ntp_interval"];
	        this.max_conn = source["max_conn"];
	        this.flow_limit = source["flow_limit"];
	        this.rate_limit = source["rate_limit"];
	        this.crypt = source["crypt"];
	        this.compress = source["compress"];
	        this.tls_enable = source["tls_enable"];
	        this.proxy_url = source["proxy_url"];
	        this.disconnect_timeout = source["disconnect_timeout"];
	        this.pprof_addr = source["pprof_addr"];
	        this.web_username = source["web_username"];
	        this.web_password = source["web_password"];
	        this.basic_username = source["basic_username"];
	        this.basic_password = source["basic_password"];
	    }
	}

}


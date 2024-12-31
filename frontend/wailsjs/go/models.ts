export namespace server {
	
	export class SetTtySizeDto {
	    tty: string;
	    col: number;
	    row: number;
	
	    static createFrom(source: any = {}) {
	        return new SetTtySizeDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tty = source["tty"];
	        this.col = source["col"];
	        this.row = source["row"];
	    }
	}
	export class GetTtyDto {
	    tty: string;
	
	    static createFrom(source: any = {}) {
	        return new GetTtyDto(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tty = source["tty"];
	    }
	}
	export class Api {
	    GetTtyDto?: GetTtyDto;
	    SetTtySizeDto?: SetTtySizeDto;
	
	    static createFrom(source: any = {}) {
	        return new Api(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.GetTtyDto = this.convertValues(source["GetTtyDto"], GetTtyDto);
	        this.SetTtySizeDto = this.convertValues(source["SetTtySizeDto"], SetTtySizeDto);
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


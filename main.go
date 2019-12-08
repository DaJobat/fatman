package main

import(
  _ "bytes"
  "encoding/binary"
  "flag"
  "fmt"
  "os"
  "time"
)

type FATEntry struct{
  name          string
  created       time.Time
  deleted       bool
}

func main(){
  var fname = flag.String("in", "", "file containing raw binary of FAT directory tree")
  flag.Parse()

  file, err := os.Open(*fname)
  if err != nil{
    panic(err)
  }

  var cur_bytes []byte = make([]byte,32)
  _, err = file.Read(cur_bytes)
  if err != nil{
    panic(err)
  }
  
  cur_entry := &FATEntry{}
  cur_entry.setFilename(cur_bytes)
  cur_entry.setTimestamps(cur_bytes)
  
  fmt.Println(cur_entry)
}

func (e *FATEntry) setFilename(bytes []byte){
  e.name = string(bytes[:11])
  e.deleted = (e.name[0] == 0xe5 || e.name[0] == 0x00)
}

func (e *FATEntry) setTimestamps(bytes []byte){
  //created_time_tenths := bytes[13]
  ctime_hms := binary.LittleEndian.Uint16(bytes[14:16])
  cdate := binary.LittleEndian.Uint16(bytes[16:18])
  
  e.created = time.Date(1980 + int(cdate >> 9),
  time.Month(int((cdate << 7) >> 12)),
  int((cdate << 11) >> 11),
  int(ctime_hms >> 11),
  int((ctime_hms << 5) >> 10),
  int((ctime_hms << 11) >> 11)*2, //2 second intervals so *2
  0,
  time.Now().Location())
//TODO: allow set local timezone
}

func (e *FATEntry) String()(string){
  return fmt.Sprintf("Name: %s\nCreated: %v\nDeleted: %t",e.name, e.created, e.deleted)
}

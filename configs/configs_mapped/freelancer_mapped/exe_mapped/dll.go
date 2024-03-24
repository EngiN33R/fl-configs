package exe_mapped

/*
This file is direct code translation from file dll_reader.py u can find in testdata folder
*/

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	gbp "github.com/darklab8/fl-configs/configs/configs_mapped/freelancer_mapped/exe_mapped/go-binary-pack"
	"golang.org/x/text/encoding/charmap"

	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind"
	"github.com/darklab8/fl-configs/configs/configs_mapped/parserutils/filefind/file"
	logus1 "github.com/darklab8/fl-configs/configs/settings/logus"
	"github.com/darklab8/go-utils/goutils/utils/utils_types"
)

type InfocardID int
type InfocardText string

const SEEK_SET = io.SeekStart // python default seek(offset, whence=os.SEEK_SET, /)

var packer = new(gbp.BinaryPack)

func Unpack[returnType any](format []string, byte_data []byte) (returnType, error) {

	unpacked_value, err := packer.UnPack(format, byte_data)
	if err != nil {
		var UnpackErrValue returnType
		return UnpackErrValue, err
	}
	value := unpacked_value[0].(returnType)
	return value, nil
}

func MakeArray(bytes_amount BytesToRead) []byte {
	switch int(bytes_amount) {
	case 1:
		return make([]byte, 1)
	case 2:
		return make([]byte, 2)
	case 4:
		return make([]byte, 4)
	case 8:
		return make([]byte, 8)
	default:
		panic("not implemented")
	}
}

func ReadUnpackWithArr[returnType any](
	fh *bytes.Reader,
	byte_data []byte,
	format []string,
) (returnType, int, error) {
	returned_n, err := fh.Read(byte_data)
	value, err := Unpack[returnType](format, byte_data)
	return value, returned_n, err
}

func ReadUnpack[returnType any](
	fh *bytes.Reader,
	bytes_amount BytesToRead,
	format []string,
) (returnType, int, error) {
	var byte_data []byte = MakeArray(bytes_amount)
	return ReadUnpackWithArr[returnType](fh, byte_data, format)
}

type BytesToRead int

func ReadUnpack2[returnType any](
	fh *bytes.Reader, bytes_amount BytesToRead,
	format []string) returnType {

	value, _, err := ReadUnpack[returnType](fh, bytes_amount, format)
	logus1.Log.CheckError(err, "failed to read unpack")
	return value
}

type DLLSection struct {
	VirtualSize          int //     DLL_Sections[name]['VirtualSize'], = struct.unpack('=l', fh.read(4))
	VirtualAddress       int //     DLL_Sections[name]['VirtualAddress'], = struct.unpack('=l', fh.read(4))
	SizeOfRawData        int //     DLL_Sections[name]['SizeOfRawData'], = struct.unpack('=l', fh.read(4))
	PointerToRawData     int //     DLL_Sections[name]['PointerToRawData'], = struct.unpack('=l', fh.read(4))
	PointerToRelocations int //     DLL_Sections[name]['PointerToRelocations'], = struct.unpack('=l', fh.read(4))
	PointerToLinenumbers int //     DLL_Sections[name]['PointerToLinenumbers'], = struct.unpack('=l', fh.read(4))
	NumberOfRelocations  int //     DLL_Sections[name]['NumberOfRelocations'], = struct.unpack('h', fh.read(2))
	NumberOfLinenumbers  int //     DLL_Sections[name]['NumberOfLinenumbers'], = struct.unpack('h', fh.read(2))
	Characteristics      int //     DLL_Sections[name]['Characteristics'], = struct.unpack('=l', fh.read(4))
}

// atatypes.append({'type': dataType, 'offset': dataOffset})
type DataType struct {
	Type_  int
	Offset int
}

var BOMcheck []byte = []byte{'\xff', '\xfe'}

func ReadText(fh *bytes.Reader, count int) string {
	strouts := [][]byte{} //     strout = b''
	total_len := 0

	for j := 0; j < count; j++ { //     for j in range(0, count):
		if j == 0 { //         if j == 0:
			h := MakeArray(2)
			fh.Read(h) //             h = fh.read(2)

			if bytes.Equal(h, BOMcheck) { //             if h == "\xff\xfe":
				continue // strip BOM
			}
			strouts = append(strouts, h) //             strout += h
			total_len += len(h)
		} else { //         else:
			portion := MakeArray(2)
			fh.Read(portion)
			strouts = append(strouts, portion) //             strout += fh.read(2)
			total_len += len(portion)
		}

	}

	result := JoinSize(total_len, strouts...)

	// PY: return strout.decode('windows-1252')[::2].encode('utf-8')
	tr := charmap.Windows1252.NewDecoder().Reader(strings.NewReader(string(result[:])))
	windows_decoded, err := io.ReadAll(tr)

	logus1.Log.CheckPanic(err, "failed to decode Windows1252")

	sliced := make([]byte, len(windows_decoded)/2)
	for i := 0; i < len(windows_decoded)/2; i += 1 {
		sliced[i] = windows_decoded[i*2] // or do whatever
	}

	return string(sliced)
}

func JoinSize(size int, s ...[]byte) []byte {
	b, i := make([]byte, size), 0
	for _, v := range s {
		i += copy(b[i:], v)
	}
	return b
}

func parseDLL(data []byte, out map[InfocardID]InfocardText, global_offset int) {
	fh := bytes.NewReader(data)

	logus1.Log.Debug("parseDLL for file.Name=")
	var returned_n64 int64
	var returned_n int
	var err error
	// Header stuff, most of it is just read and ignored but we need a few addresses from it.

	returned_n64, err = fh.Seek(0x3C, SEEK_SET) // fh.seek(0x3C)
	PE_sig_loc, returned_n, err := ReadUnpack[int](fh, BytesToRead(1), []string{"B"})

	returned_n64, err = fh.Seek(int64(PE_sig_loc+4), SEEK_SET)                                        // fh.seek(PE_sig_loc + 4) # goto COFF header (after sig)
	returned_n, err = fh.Read(make([]byte, 2))                                                        // COFF_Head_Machine, = struct.unpack('h', fh.read(2)) # 014c - i386 or compatible
	COFF_Head_NumberOfSections, returned_n, err := ReadUnpack[int](fh, BytesToRead(2), []string{"h"}) // COFF_Head_NumberOfSections, = struct.unpack('h', fh.read(2))
	returned_n, err = fh.Read(make([]byte, 4))                                                        // COFF_Head_TimeDateStamp, = struct.unpack('=l', fh.read(4))
	returned_n, err = fh.Read(make([]byte, 4))                                                        // COFF_Head_PointerToSymbolTable, = struct.unpack('=l', fh.read(4))
	returned_n, err = fh.Read(make([]byte, 4))                                                        // COFF_Head_NumberOfSymbols, = struct.unpack('=l', fh.read(4))

	COFF_Head_SizeOfOptionalHeader, returned_n, err := ReadUnpack[int](fh, BytesToRead(2), []string{"h"}) // COFF_Head_SizeOfOptionalHeader, = struct.unpack('h', fh.read(2))
	COFF_Head_Characteristics, _, err := ReadUnpack[int](fh, BytesToRead(2), []string{"h"})               // COFF_Head_Characteristics, = struct.unpack('h', fh.read(2)) # 210e
	_ = COFF_Head_Characteristics

	OPT_Head_Start, err := fh.Seek(0, io.SeekCurrent)

	if COFF_Head_SizeOfOptionalHeader != 0 { // if COFF_Head_SizeOfOptionalHeader != 0: # image header exists
		fh.Read(make([]byte, 2)) //     OPT_Head_Magic, = struct.unpack('h', fh.read(2))
		fh.Read(make([]byte, 1)) //     OPT_Head_MajorLinkerVers, = struct.unpack('c', fh.read(1))
		fh.Read(make([]byte, 1)) //     OPT_Head_MinorLinkerVers, = struct.unpack('c', fh.read(1))
		fh.Read(make([]byte, 4)) //     OPT_Head_SizeOfCode, = struct.unpack('=l', fh.read(4))
		fh.Read(make([]byte, 4)) //     OPT_Head_SizeOfInitializedData, = struct.unpack('=l', fh.read(4))
		fh.Read(make([]byte, 4)) //     OPT_Head_SizeOfUninitializedData, = struct.unpack('=l', fh.read(4))
		fh.Read(make([]byte, 4)) //     OPT_Head_AddressOfEntryPoint, = struct.unpack('=l', fh.read(4))
		fh.Read(make([]byte, 4)) //     OPT_Head_BaseOfCode, = struct.unpack('=l', fh.read(4))

		//     if OPT_Head_Magic == 0x20B: # if it's 64-bit
		//         OPT_Head_ImageBase, = struct.unpack('q', fh.read(8))
		//     else:
		//         OPT_Head_BaseOfData, = struct.unpack('=l', fh.read(4))
		//         OPT_Head_ImageBase, = struct.unpack('=l', fh.read(4))
		fh.Read(make([]byte, 8))

		fh.Read(make([]byte, 4)) //     SectionAlignment = fh.read(4)
		fh.Read(make([]byte, 4)) //     FileAlignment = fh.read(4)
		fh.Read(make([]byte, 2)) //     MajorOperatingSystemVersion = fh.read(2)
		fh.Read(make([]byte, 2)) //     MinorOperatingSystemVersion = fh.read(2)
		fh.Read(make([]byte, 2)) //     MajorImageVersion = fh.read(2)
		fh.Read(make([]byte, 2)) //     MinorImageVersion = fh.read(2)
		fh.Read(make([]byte, 2)) //     MajorSubsystemVersion = fh.read(2)
		fh.Read(make([]byte, 2)) //     MinorSubsystemVersion = fh.read(2)
		fh.Read(make([]byte, 4)) //     Win32VersionValue = fh.read(4)
		fh.Read(make([]byte, 4)) //     SizeOfImage = fh.read(4)
		fh.Read(make([]byte, 4)) //     SizeOfHeaders = fh.read(4)

	}

	// # Get the section header info, we only care about ".rsrc" though
	fh.Seek(int64(OPT_Head_Start)+int64(COFF_Head_SizeOfOptionalHeader), SEEK_SET) // fh.seek(OPT_Head_Start + COFF_Head_SizeOfOptionalHeader)
	var DLL_Sections map[string]*DLLSection = make(map[string]*DLLSection)         // DLL_Sections = {}
	for i := 0; i < int(COFF_Head_NumberOfSections); i++ {                         // for i in range(0, COFF_Head_NumberOfSections):
		logus1.Log.Debug("i := 0; i < int(COFF_Head_NumberOfSections); i++, i=" + strconv.Itoa(i))
		//     nt = fh.read(8)
		nt := make([]byte, 8)
		fh.Read(nt)

		name := strings.ReplaceAll(string(nt), "\x00", "") //     name = nt.decode('utf-8').strip("\x00") # TODO: There was much more complex code for this in PHP, but the input format is completely different. Like different order and format different.

		section := &DLLSection{}
		DLL_Sections[name] = section
		section.VirtualSize = ReadUnpack2[int](fh, BytesToRead(4), []string{"l"})          // LL_Sections[name]['VirtualSize'], = struct.unpack('=l', fh.read(4))
		section.VirtualAddress = ReadUnpack2[int](fh, BytesToRead(4), []string{"l"})       //     DLL_Sections[name]['VirtualAddress'], = struct.unpack('=l', fh.read(4))
		section.SizeOfRawData = ReadUnpack2[int](fh, BytesToRead(4), []string{"l"})        //     DLL_Sections[name]['SizeOfRawData'], = struct.unpack('=l', fh.read(4))
		section.PointerToRawData = ReadUnpack2[int](fh, BytesToRead(4), []string{"l"})     //     DLL_Sections[name]['PointerToRawData'], = struct.unpack('=l', fh.read(4))
		section.PointerToRelocations = ReadUnpack2[int](fh, BytesToRead(4), []string{"l"}) //     DLL_Sections[name]['PointerToRelocations'], = struct.unpack('=l', fh.read(4))
		section.PointerToLinenumbers = ReadUnpack2[int](fh, BytesToRead(4), []string{"l"}) //     DLL_Sections[name]['PointerToLinenumbers'], = struct.unpack('=l', fh.read(4))
		section.NumberOfRelocations = ReadUnpack2[int](fh, BytesToRead(2), []string{"h"})  //     DLL_Sections[name]['NumberOfRelocations'], = struct.unpack('h', fh.read(2))
		section.NumberOfLinenumbers = ReadUnpack2[int](fh, BytesToRead(2), []string{"h"})  //     DLL_Sections[name]['NumberOfLinenumbers'], = struct.unpack('h', fh.read(2))
		section.Characteristics = ReadUnpack2[int](fh, BytesToRead(4), []string{"l"})      //     DLL_Sections[name]['Characteristics'], = struct.unpack('=l', fh.read(4))

	}

	logus1.Log.Debug("rsrcstart")
	rsrcstart := DLL_Sections[".rsrc"].PointerToRawData               // rsrcstart = DLL_Sections['.rsrc']['PointerToRawData']
	fh.Seek(int64(rsrcstart)+int64(14), io.SeekStart)                 // fh.seek(rsrcstart + 14) # go to start of .rsrc
	numentries := ReadUnpack2[int](fh, BytesToRead(2), []string{"h"}) // numentries, = struct.unpack('h', fh.read(2))
	datatypes := []*DataType{}
	// # get the data types stored in the resource section
	for i := 0; i < numentries; i++ { // for i in range(0, numentries):
		logus1.Log.Debug("for i := 0; i < numentries; i++, i=" + strconv.Itoa(i))

		dataType := ReadUnpack2[int](fh, BytesToRead(4), []string{"l"}) //     dataType, = struct.unpack('=l', fh.read(4))

		doi := make([]byte, 2)
		fh.Read(doi) //     doi = fh.read(2)
		doj := make([]byte, 1)
		fh.Read(doj) //     doj = fh.read(1)

		//     dataOffset, = struct.unpack('<i', doi + doj + '\x00'.encode('utf-8'))
		packer := new(gbp.BinaryPack)
		unpacked_value, err := packer.UnPack([]string{"i"}, bytes.Join([][]byte{doi, doj, []byte{'\x00'}}, []byte{}))
		logus1.Log.CheckError(err, "failed to unpack")
		dataOffset := unpacked_value[0].(int)

		datatypes = append(datatypes, &DataType{
			Type_:  dataType,
			Offset: dataOffset,
		}) //     datatypes.append({'type': dataType, 'offset': dataOffset})
		fh.Seek(1, io.SeekCurrent) //     fh.seek(1, os.SEEK_CUR) # jump ahead 1 byte
	}

	// # each different data type is stored in a block, loop through each
	for _, datatype := range datatypes { // for i in range(0, len(datatypes)):
		logus1.Log.Debug("for _, datatype := range datatypes {" + fmt.Sprintf("%v", datatype))
		fh.Seek(int64(datatype.Offset)+int64(rsrcstart), io.SeekStart) //     fh.seek(datatypes[i]['offset'] + rsrcstart)

		name := MakeArray(8)
		fh.Read(name) //     name = fh.read(8)

		fh.Seek(6, io.SeekCurrent)                                        //     fh.seek(6, os.SEEK_CUR)
		numentries := ReadUnpack2[int](fh, BytesToRead(2), []string{"h"}) //     numentries, = struct.unpack('h', fh.read(2))

		fh.Seek(0, io.SeekCurrent) //     sectionstart = fh.tell() # remember where we are here

		for entry := 0; entry < numentries; entry++ { // for entry in range(0, numentries):                   //     for entry in range(0, numentries):
			logus1.Log.Debug("for entry := 0; entry < numentries; entry++ entry=" + strconv.Itoa(entry))
			//         # get the id number and location of this entry
			idnum := ReadUnpack2[int](fh, BytesToRead(4), []string{"i"}) //         idnum, = struct.unpack('i', fh.read(4))

			doi := MakeArray(2)
			fh.Read(doi) //     doi = fh.read(2)
			doj := MakeArray(1)
			fh.Read(doj) //     doj = fh.read(1)

			//         nameloc, = struct.unpack('<i', doi + doj + '\x00'.encode('utf-8'))
			packer := new(gbp.BinaryPack)
			unpacked_value, err := packer.UnPack([]string{"i"}, JoinSize(len(doi)+len(doj)+1, doi, doj, []byte{'\x00'}))
			logus1.Log.CheckError(err, "failed to unpack")
			nameloc := unpacked_value[0].(int)

			brk := MakeArray(1)
			fh.Read(brk) //         brk = fh.read(1)

			backto, _ := fh.Seek(0, io.SeekCurrent) //         backto = fh.tell() # remember where we were in the list of entries

			fh.Seek(int64(rsrcstart)+int64(nameloc), io.SeekStart) //         fh.seek(rsrcstart + nameloc) # jump to the entry

			name := MakeArray(8)
			fh.Read(name)              //         name = fh.read(8) # get the name
			fh.Seek(8, io.SeekCurrent) //         fh.seek(8, os.SEEK_CUR)

			lang := MakeArray(4)
			fh.Read(lang) //         lang = fh.read(4) # language for this resource

			someinfoloc := ReadUnpack2[int](fh, BytesToRead(4), []string{"i"}) //         someinfoloc, = struct.unpack('i', fh.read(4)) # location of the real location of the entry....

			fh.Seek(int64(rsrcstart)+int64(someinfoloc), SEEK_SET)            //         fh.seek(rsrcstart + someinfoloc) # jump there
			absloc := ReadUnpack2[int](fh, BytesToRead(4), []string{"i"})     //         absloc, = struct.unpack('i', fh.read(4)) # get the real location
			datalength := ReadUnpack2[int](fh, BytesToRead(4), []string{"i"}) //         datalength, = struct.unpack('i', fh.read(4)) # entry length in bytes

			GetResource(data, out, absloc, datatype, idnum, global_offset, datalength)

			//         # go back and get the next one
			fh.Seek(backto, io.SeekStart) //         fh.seek(backto)
		}
	}

	_ = returned_n
	_ = returned_n64
	_ = err
}

func GetResource(
	data []byte,
	out map[InfocardID]InfocardText,
	absloc int,
	datatype *DataType,
	idnum int,
	global_offset int,
	datalength int,
) error {
	fh := bytes.NewReader(data)

	//         # now that we've got absolute location of each resource, get it!
	fh.Seek(int64(absloc), io.SeekStart) //         fh.seek(absloc)

	if datatype.Type_ == 0x06 { //         if datatypes[i]['type'] == 0x06: # string table
		for strindex := 0; strindex < 16; strindex++ { //             for strindex in range(0, 16): # each string table has up to 16 entries
			tableLen, n, err := ReadUnpack[int](fh, BytesToRead(2), []string{"h"}) //                 tableLen, = struct.unpack('h', fh.read(2))
			//                 if not tableLen:
			//                     continue # drop completely empty strings
			if tableLen == 0 || n == 0 || err != nil {
				continue
			}

			ids_index := (idnum-1)*16 + strindex + global_offset //                 ids_index = (idnum - 1)*16 + strindex + global_offset
			ids_text := ReadText(fh, tableLen)                   //                 ids_text = ReadText(fh, tableLen)

			out[InfocardID(ids_index)] = InfocardText(ids_text) //                 out[ids_index] = ids_text
		}

	} else if datatype.Type_ == 0x17 { //         elif datatypes[i]['type'] == 0x17: # html
		ids_index := idnum + global_offset //             ids_index = idnum + global_offset
		if datalength%2 != 0 {             //             if datalength % 2:
			datalength -= 1 //                 datalength -= 1 # if odd length, ignore the last byte (UTF-16 is 2 bytes per character...)
		}

		if 500904 == ids_index {
			_ = datalength
		}

		floored_datalength := math.Floor(float64(datalength) / 2)
		ids_text := ReadText(fh, int(floored_datalength)) //             ids_text = ReadText(fh, datalength // 2).rstrip()

		out[InfocardID(ids_index)] = InfocardText(ids_text) //             out[ids_index] = ids_text
	}
	return nil
}

func ParseDLLs(dll_fnames []*file.File) map[InfocardID]InfocardText {
	out := make(map[InfocardID]InfocardText, 0)

	for idx, name := range dll_fnames {
		data, err := os.ReadFile(name.GetFilepath().ToString())

		if logus1.Log.CheckError(err, "unable to read dll") {
			continue
		}

		global_offset := int(math.Pow(2, 16)) * (idx + 1)
		parseDLL(data, out, global_offset)
	}

	return out
}

func GetAllInfocards(filesystem *filefind.Filesystem, dll_names []string) map[InfocardID]InfocardText {

	var files []*file.File
	for _, filename := range dll_names {
		dll_file := filesystem.GetFile(utils_types.FilePath(strings.ToLower(filename)))
		files = append(files, dll_file)
	}

	return ParseDLLs(files)
}

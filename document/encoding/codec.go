package encoding

import (
	"io"
	"tgenj/document"
)

/**
 *  @ClassName:codec
 *  @Description:TODO
 *  @Author:jackey
 *  @Create:2021/3/11 下午9:34
 */

// A Codec is able to create encoders and decoders for a specific encoding format.
type Codec interface {
	NewEncoder(io.Writer) Encoder
	// NewDocument returns a document without decoding its given binary representation.
	// The returned document should ideally support random-access, i.e. decoding one path
	// without decoding the entire document. If not, the document must be lazily decoded.
	NewDocument([]byte) document.Document
}

// An Encoder encodes one document to the underlying writer.
type Encoder interface {
	EncodeDocument(d document.Document) error
	// Close the encoder to release any resource.
	Close()
}

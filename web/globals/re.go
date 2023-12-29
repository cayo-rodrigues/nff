package globals

import "regexp"

var ReIeMg = regexp.MustCompile(`^\d{3}.?\d{3}.?\d{3}\/?\d{4}$`)
var ReCpf = regexp.MustCompile(`^\d{3}.?\d{3}.?\d{3}\-?\d{2}$`)
var ReCnpj = regexp.MustCompile(`^(\d{2}.?\d{3}.?\d{3}\/?\d{4}\-?\d{2})$`)
var ReEmail = regexp.MustCompile(`[^@ \t\r\n]+@[^@ \t\r\n]+\.[^@ \t\r\n]+`)
var RePostalCode = regexp.MustCompile(`(^\d{5})\-?(\d{3}$)`)
var ReAddressNumber = regexp.MustCompile(`^(?:s\/n|S\/n|S\/N|s\/N)|^(\d)*$`)
var ReSiareNFANumber = regexp.MustCompile(`^(\d{3}\.\d{3}\.\d{3}|\d{9})$`)
var ReSiareNFAProtocol = regexp.MustCompile(`^(\d{13})$`)
var ReGta = regexp.MustCompile(`^([a-zA-Z]-\d{1,6}(;\s*[a-zA-Z]-\d{1,6})*$)`)

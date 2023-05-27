export function createInvoicesPage() {
    document.querySelector("#current-tab-title").innerText = "Notas Fiscais"

    const className = "sub-menu__item--selected"

    const subEntry1 = document.querySelector("#sub-entry-1")
    subEntry1.innerText = "Emitir Notas Fiscais"
    subEntry1.classList.add(className)

    const subEntry2 = document.querySelector("#sub-entry-2")
    subEntry2.innerText = "Cancelar Notas Fiscais"
    subEntry2.classList.remove(className)

    const contentCore = document.querySelector("#content__core")
    contentCore.innerHTML = ""
    contentCore.innerHTML = `
        <form class="invoices-form">
            <div class="invoices-form__inputs-container">

                <div class="invoices-form__input">
                    <label for="sender-input">Remetente</label>
                    <select name="sender" id="sender-input">
                        <option value="sender-id">Emerson</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="recipient-input">Destinatário</label>
                    <select name="recipient" id="recipient-input">
                        <option value="emerson-id">Emerson</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="operation-input">Natureza da Operação</label>
                    <select name="operation" id="operation-input">
                        <option value="VENDA">VENDA</option>
                        <option value="REMESSA">REMESSA</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="gta-input">GTA</label>
                    <input type="gta" name="gta" id="gta-input">
                </div>
                <div class="invoices-form__input">
                    <label for="cfop-input">CFOP</label>
                    <select name="cfop" id="cfop-input">
                        <option value="5101">5101</option>
                        <option value="5102">5102</option>
                        <option value="5103">5103</option>
                        <option value="5105">5105</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="shipping-input">Frete</label>
                    <input type="text" name="shipping" id="shipping-input">
                </div>
                <div class="invoices-form__input">
                    <label for="add_shipping_to_total_value-input">Adicionar Frete ao Total</label>
                    <select name="add_shipping_to_total_value" id="add_shipping_to_total_value-input">
                        <option value="sim">Sim</option>
                        <option value="não">Não</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="is_final_customer-input">Consumidor Final</label>
                    <select name="is_final_customer" id="is_final_customer-input">
                        <option value="sim">Sim</option>
                        <option value="não">Não</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="icms-input">Contribuinte ICMS</label>
                    <select name="icms" id="icms-input">
                        <option value="sim">Sim</option>
                        <option value="não">Não</option>
                    </select>
                </div>
                <div class="invoices-form__input">
                    <label for="custom_file_name-input">Nome do Arquivo</label>
                    <input type="text" name="custom_file_name" id="custom_file_name-input">
                </div>
                <div class="invoices-form__input">
                    <label for="extra_notes-input">Informações Complementares</label>
                    <input type="text" name="extra_notes" id="extra_notes-input">
                </div>
                <div class="invoices-form__input">
                    <label>Itens da Nota Fiscal</label>
                    <button type="button" class="invoices-form__itens-modal-button">Gerenciar Itens</button>
                </div>

            </div>

            <button class="invoices-form__button" type="submit">Emitir Notas Fiscais</button>
        </form>
    `
}

export function cancelInvoicesPage() {
    const contentCore = document.querySelector("#content__core")
    contentCore.innerHTML = ""
    contentCore.innerHTML = `
        Cancel invoices!
    `
}

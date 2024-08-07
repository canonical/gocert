import { Dispatch, SetStateAction, useState, ChangeEvent } from "react"
import { useMutation, useQueryClient } from "react-query"
import { ConfirmationModalData } from "./row"
import { extractCert, csrMatchesCertificate } from "../utils"
import { postCertToID } from "../queries"
import { useCookies } from "react-cookie"

interface ConfirmationModalProps {
    modalData: ConfirmationModalData
    setModalData: Dispatch<SetStateAction<ConfirmationModalData>>
}


export function ConfirmationModal({ modalData, setModalData }: ConfirmationModalProps) {
    const confirmQuery = () => {
        modalData?.onMouseDownFunc()
        setModalData(null)
    }
    return (
        <div className="p-modal" id="modal">
            <section className="p-modal__dialog" role="dialog" aria-modal="true" aria-labelledby="modal-title" aria-describedby="modal-description">
                <header className="p-modal__header">
                    <h2 className="p-modal__title" id="modal-title">Confirm Action</h2>
                </header>
                <p>{modalData?.warningText}</p>
                <footer className="p-modal__footer">
                    <button className="u-no-margin--bottom" aria-controls="modal" onMouseDown={() => setModalData(null)}>Cancel</button>
                    <button className="p-button--negative u-no-margin--bottom" onMouseDown={confirmQuery}>Confirm</button>
                </footer>
            </section>
        </div>
    )
}

function SubmitCertificate({ existingCSRText, existingCertText, certText, onClickFunc }: { existingCSRText: string, existingCertText: string, certText: string, onClickFunc: any }) {
    let certIsValid = false
    let certMatchesCSR = false
    try {
        extractCert(certText)
        certIsValid = true
        if (csrMatchesCertificate(existingCSRText, certText)) {
            certMatchesCSR = true
        }
    }
    catch { }

    const validationComponent = certText == "" ?
        <></> :
        !certIsValid ?
            <div><i className="p-icon--error"></i> Invalid Certificate</div> :
            existingCertText == certText ?
                <div><i className="p-icon--error"></i> Certificate is identical to the one uploaded</div> :
                !certMatchesCSR ?
                    <div><i className="p-icon--error"></i> Certificate does not match the request</div> :
                    <div><i className="p-icon--success"></i> Valid Certificate</div>
    const buttonComponent = certIsValid && certMatchesCSR && existingCertText != certText ? <button className="p-button--positive" name="submit" onClick={onClickFunc} >Submit</button> : <button className="p-button--positive" name="submit" disabled={true} onClick={onClickFunc} >Submit</button>
    return (
        <>
            {validationComponent}
            {buttonComponent}
        </>
    )
}

interface SubmitCertificateModalProps {
    id: string
    csr: string
    cert: string
    setFormOpen: Dispatch<SetStateAction<boolean>>
}
export function SubmitCertificateModal({ id, csr, cert, setFormOpen }: SubmitCertificateModalProps) {
    const [cookies, setCookie, removeCookie] = useCookies(['user_token']);
    const queryClient = useQueryClient()
    const mutation = useMutation(postCertToID, {
        onSuccess: () => {
            queryClient.invalidateQueries('csrs')
        },
    })
    const [certificatePEMString, setCertificatePEMString] = useState<string>("")
    const handleTextChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
        setCertificatePEMString(event.target.value);
    }
    const handleFileChange = (event: ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0]
        if (file) {
            const reader = new FileReader();
            reader.onload = (e: ProgressEvent<FileReader>) => {
                if (e.target) {
                    if (e.target.result) {
                        setCertificatePEMString(e.target.result.toString());
                    }
                }
            };
            reader.readAsText(file);
        }
    };
    const handleSubmit = () => {
        mutation.mutate({ id: id, authToken: cookies.user_token, cert: certificatePEMString })
        setFormOpen(false)
    }
    return (
        <div className="p-modal" id="modal">
            <section className="p-modal__dialog" role="dialog" aria-modal="true" aria-labelledby="modal-title" aria-describedby="modal-description">
                <header className="p-modal__header">
                    <h2 className="p-modal__title" id="modal-title">Submit Certificate</h2>
                </header>
                <form className="p-form p-form--stacked">
                    <div className="p-form__group row">
                        <label htmlFor="textarea">
                            Enter or upload the Certificate in PEM format below
                        </label>
                        <textarea id="csr-textarea" name="textarea" rows={10} placeholder="-----BEGIN CERTIFICATE-----" onChange={handleTextChange} value={certificatePEMString} />
                    </div>
                    <div className="p-form__group row">
                        <input type="file" name="upload" accept=".pem,.crt" onChange={handleFileChange}></input>
                    </div>
                    <div className="p-form__group row">
                    </div>
                </form>
                <footer className="p-modal__footer">
                    <SubmitCertificate existingCSRText={csr.trim()} existingCertText={cert.trim()} certText={certificatePEMString.trim()} onClickFunc={handleSubmit} />
                    <button className="u-no-margin--bottom" aria-controls="modal" onMouseDown={() => setFormOpen(false)}>Cancel</button>
                </footer>
            </section>
        </div>
    )
}

export function SuccessNotification({ successMessage }: { successMessage: string }) {
    const style = {
        display: 'inline'
    }
    return (
        <p style={style}>
            <i className="p-icon--success"></i> {successMessage}
        </p>
    );
}
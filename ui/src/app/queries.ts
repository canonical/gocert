import { CSREntry } from "./types"

export async function getCertificateRequests(): Promise<CSREntry[]> {
    const response = await fetch("/api/v1/certificate_requests")
    if (!response.ok) {
        throw new Error('Network response was not ok')
    }
    return response.json()
}

export async function postCSR(csr: string) {
    const response = await fetch("/api/v1/certificate_requests", {
        method: 'post',
        headers: {
            'Content-Type': 'text/plain',
        },
        body: csr.trim()
    })
    if (!response.ok) {
        throw new Error('Network response was not ok')
    }
    return response.json()
}
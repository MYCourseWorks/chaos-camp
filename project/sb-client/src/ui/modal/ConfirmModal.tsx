import React from "react"

interface Props {
    isOpen: boolean
    title: string
    onSubmit(): void
    onCancel(): void
}

export class ConfirmModal extends React.Component<Props> {

    public render(): React.ReactNode {
        const { children, isOpen, title } = this.props

        return (
            <div className="modal" style={{ display: isOpen ? "block" : "none"}}>
                <div className="modal-background" />
                <div className="modal-card">
                    <header className="modal-card-head">
                        <p className="modal-card-title">{title}</p>
                        <button className="delete" onClick={this.props.onCancel} />
                    </header>
                    <section className="modal-card-body">
                        {children}
                    </section>
                    <footer className="modal-card-foot">
                        <button className="button" onClick={this.props.onSubmit}>Ok</button>
                        <button className="button" onClick={this.props.onCancel}>Cancel</button>
                    </footer>
                </div>
            </div>
        )
    }
}
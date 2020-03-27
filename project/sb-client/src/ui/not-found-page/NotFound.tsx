import React from "react"
import "./NotFound.css"

export class NotFound extends React.Component {
    public render(): React.ReactNode {
        return (
            <div id="fof-main">
                <div className="fof">
                    <h1>Error 404, Page Not Found!</h1>
                </div>
            </div>
        )
    }
}
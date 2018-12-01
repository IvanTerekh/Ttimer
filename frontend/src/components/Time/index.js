import React from "react";
import format from '../../logic/format';
import './Time.css';

function currentDatetime() {
    let d = new Date();
    let s = d.toISOString();
    return s.replace('T', ' ').slice(0, -2);
}

export default class Time extends React.Component {

    constructor(props) {
        super(props);

        this.state = {
            timeState: "usual",
            centis: 0
        };

        this.running = false;
        this.keydown = this.keydown.bind(this);
        this.keyup = this.keyup.bind(this);
        this.start = this.start.bind(this);
        this.stop = this.stop.bind(this);
        document.onkeydown = (e) => {
            if (e.code === "Space") this.keydown(e)
        };
        document.onkeyup = (e) => {
            if (e.code === "Space") this.keyup(e)
        };
    }

    start() {
        this.running = true;
        this.startTime = Date.now();

        this.timerId = setInterval(
            () => this.tick(), 1000 / 60);
        document.ontouchstart = this.keydown;
        document.onkeydown = this.keydown;
    }

    stop() {
        const res = Math.floor((Date.now() - this.startTime) / 10);
        clearInterval(this.timerId);
        this.props.saveResult({
            centis: res,
            scramble: this.props.scramble,
            penalty: false,
            datetime: currentDatetime()
        });
        this.running = false;
        this.setState({
            centis: res
        });
        document.ontouchstart = () => {
        };
        document.onkeydown = (e) => {
            if (e.code === "Space") this.keydown(e)
        };
    }

    keydown(e) {
        e.preventDefault();
        if (this.running) {
            this.stop();
        } else if (isNaN(this.timeoutId)) {
            this.setState({
                timeState: "notReady"
            });
            this.timeoutId = setTimeout(() => {
                this.setState({
                    timeState: "ready"
                });
            }, 300);
        }
    }

    keyup(e) {
        if (!this.running) {
            if (this.state.timeState === "ready") {
                this.start()
            }
            clearTimeout(this.timeoutId);
            this.timeoutId = NaN;

            this.setState({
                timeState: "usual"
            });
        }
    }

    tick() {
        this.setState({
            centis: Math.floor((Date.now() - this.startTime) / 10)
        });
    }

    render() {
        return <div onTouchStart={this.keydown} title='Hold "Space" to start.' onTouchEnd={this.keyup}>
            <h1 className={"time " + this.state.timeState}>{format(this.state.centis)}</h1>
        </div>
    }
}

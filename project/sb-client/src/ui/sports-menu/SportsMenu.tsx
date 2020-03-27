import React from "react"
import "./SportsMenu.css"
import { MainStore } from "../../stores/mainStore"
import { League } from "../../models/League"
import { Link } from "react-router-dom";
import DownIcon from "../../icons/down.svg";
import RightIcon from "../../icons/right.svg";

interface Props {
    mainStore: MainStore
}

interface SportView {
    id: number,
    name: string,
    isSelected: boolean,
    countries: Map<string, CountryView>,
}

interface CountryView {
    name: string,
    isSelected: boolean,
    leagues: Map<number, LeagueView>,
}

interface LeagueView {
    id: number,
    name: string,
}

interface State {
    sports: Map<number, SportView>
}

export class SportsMenu extends React.Component<Props, State> {

    constructor(props: Props) {
        super(props);
        this.state = { sports: new Map() }
    }

    public async componentDidMount(): Promise<void> {
        try {
            const leagues = await this.props.mainStore.services.leagues.all()
            const sportViews: Map<number, SportView> = new Map()

            leagues.forEach((l: League) => {
                let s: SportView
                if (sportViews.has(l.sport.id)) {
                    s = sportViews.get(l.sport.id) as SportView
                } else {
                    s = {
                        id: l.sport.id,
                        isSelected: true,
                        name: l.sport.name,
                        countries: new Map()
                    }
                    sportViews.set(s.id, s)
                }

                let c: CountryView
                if (s.countries.has(l.country)) {
                    c = s.countries.get(l.country) as CountryView
                } else {
                    c = { name: l.country, leagues: new Map(), isSelected: false }
                    s.countries.set(c.name, c)
                }

                let lv: LeagueView
                if (c.leagues.has(l.id)) {
                    lv = c.leagues.get(l.id) as LeagueView
                } else {
                    lv = { id: l.id, name: l.name }
                    c.leagues.set(lv.id, lv)
                }
            })

            this.setState({ sports: sportViews })
        } catch (err) {
            console.error(err)
        }
    }

    public render(): React.ReactNode {
        return (
            <div className="sports-menu">
                <ul className="sport-menu-sports">
                    {this.renderSports()}
                </ul>
            </div>
        );
    }

    public renderSports(): React.ReactNode[] {
        const ret: React.ReactNode[] = []
        this.state.sports.forEach((s: SportView) =>
            ret.push(
                <React.Fragment key={s.id}>
                    <li
                        className="sport-menu-sport-item"
                        onClick={() => {
                            s.isSelected = !s.isSelected
                            this.setState({ sports: this.state.sports })
                        }}
                    >
                        <span className="sport-menu-icon-span">
                            {
                                s.isSelected ?
                                    <img src={DownIcon} alt="down icon" width={15} height={15}/> :
                                    <img src={RightIcon} alt="down icon" width={10} height={10}/>
                            }
                        </span>
                        {s.name}
                    </li>
                    {
                        s.isSelected &&
                        s.countries?.size > 0 &&
                        (<ul className="sports-menu-countries">
                            {this.renderCountries(s)}
                        </ul>)
                    }
                </React.Fragment>
            )
        )

        return ret
    }

    public renderCountries(s: SportView): React.ReactNode[] {
        const ret: React.ReactNode[] = []
        s.countries?.forEach((c: CountryView) => {
            ret.push(
                <React.Fragment key={c.name}>
                    <li
                        className="sports-menu-country-item"
                        onClick={() => {
                            c.isSelected = !c.isSelected
                            this.setState({ sports: this.state.sports })
                        }}
                    >
                        <span className="sport-menu-icon-span">
                            {
                                c.isSelected ?
                                    <img src={DownIcon} alt="down icon" width={15} height={15}/> :
                                    <img src={RightIcon} alt="down icon" width={10} height={10}/>
                            }
                        </span>
                        {c.name}
                    </li>
                    {
                        c.isSelected &&
                        c.leagues?.size > 0 &&
                        (<ul className="sports-menu-leagues">
                            {this.renderLeagues(c)}
                        </ul>)
                    }
                </React.Fragment>
            )
        })

        return ret
    }

    public renderLeagues(c: CountryView): React.ReactNode[] {
        const ret: React.ReactNode[] = []
        c.leagues?.forEach((l: LeagueView) => {
            ret.push(
                <li className="sports-menu-league-item" key={l.id}>
                    <Link to={`/leagues/${l.id}`}>
                        {l.name}
                    </Link>
                </li>
            )
        })

        return ret
    }
}
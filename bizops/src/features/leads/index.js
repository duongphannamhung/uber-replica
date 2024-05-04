import moment from "moment"
import { useEffect } from "react"
import { useDispatch, useSelector } from "react-redux"
import TitleCard from "../../components/Cards/TitleCard"
import { openModal } from "../common/modalSlice"
import { deleteLead, getTrips } from "./leadSlice"
import { CONFIRMATION_MODAL_CLOSE_TYPES, MODAL_BODY_TYPES } from '../../utils/globalConstantUtil'
import TrashIcon from '@heroicons/react/24/outline/TrashIcon'
import { showNotification } from '../common/headerSlice'
import { useState } from "react"

const TopSideButtons = () => {

    const dispatch = useDispatch()

    const openAddNewLeadModal = () => {
        dispatch(openModal({title : "Add New Lead", bodyType : MODAL_BODY_TYPES.LEAD_ADD_NEW}))
    }

    return(
        <div className="inline-block float-right">
            <button className="btn px-6 btn-sm normal-case btn-primary" onClick={() => openAddNewLeadModal()}>Book trip</button>
        </div>
    )
}

function Leads(){

    const {leads } = useSelector(state => state.lead)
    const dispatch = useDispatch()
    const [page, setPage] = useState(1);

    useEffect(() => {
        dispatch(getTrips())
    }, [])

    

    const getDummyStatus = (index) => {
        if(index % 5 === 0)return <div className="badge">Not Interested</div>
        else if(index % 5 === 1)return <div className="badge badge-primary">In Progress</div>
        else if(index % 5 === 2)return <div className="badge badge-secondary">Sold</div>
        else if(index % 5 === 3)return <div className="badge badge-accent">Need Followup</div>
        else return <div className="badge badge-ghost">Open</div>
    }

    const getVehicleType = (index) => {
        if(index === 1)return <div className="badge badge-secondary">UrepBike</div>
        else if(index === 2)return <div className="badge badge-primary">UrepCar</div>
        else if(index === 3)return <div className="badge badge-accent">UrepCar7</div>
        else return <div className="badge badge-info">UrepPlus</div>
    }

    const getIsStarted = (is_started) => {
        if(is_started) return <div className="badge badge-success">Yes</div>
        else return <div className="badge badge-ghost">No</div>
    }

    const deleteCurrentLead = (index) => {
        dispatch(openModal({title : "Confirmation", bodyType : MODAL_BODY_TYPES.CONFIRMATION, 
        extraObject : { message : `Are you sure you want to delete this lead?`, type : CONFIRMATION_MODAL_CLOSE_TYPES.LEAD_DELETE, index}}))
    }

    return(
        <>
            <TitleCard title="Current Trips" topMargin="mt-2" TopSideButtons={<TopSideButtons />}>
                {/* Leads List in table format loaded from slice after api call */}
            <div className="overflow-x-auto w-full">
                <table className="table w-full">
                    <thead>
                    <tr>
                        <th>Trip Id</th>
                        <th>Driver Id</th>
                        <th>Service Type</th>
                        <th>Is Started</th>
                        <th>Departure Place</th>
                        <th>Destination Place</th>
                        <th>Fare</th>
                        <th>Time Created</th>
                        <th></th>
                    </tr>
                    </thead>
                    <tbody>
                        {
                            leads.map((l, k) => {
                                return(
                                    <tr key={k}>
                                    <td style={{ textAlign: 'center' }}>{l.id}</td>
                                    <td style={{ textAlign: 'center' }}>{l.driver_id.Int32}</td>
                                    <td style={{ textAlign: 'center' }}>{getVehicleType(l.service_type)}</td>
                                    <td style={{ textAlign: 'center' }}>{getIsStarted(l.is_started)}</td>
                                    <td style={{ textAlign: 'center' }}>{l.departure_name}</td>
                                    <td style={{ textAlign: 'center' }}>{l.destination_name}</td>
                                    <td style={{ textAlign: 'center' }}>{l.fare.Int32}</td>
                                    <td style={{ textAlign: 'center' }}>{l.created_at}</td>
                                    <td><button className="btn btn-square btn-ghost" onClick={() => deleteCurrentLead(k)}><TrashIcon className="w-5"/></button></td>
                                    </tr>
                                )
                            })
                        }
                    </tbody>
                </table>
            </div>
            </TitleCard>
        </>
    )
}

{/* <td>{moment(new Date()).add(-5*(k+2), 'days').format("DD MMM YY")}</td> 
            <button disabled={page === 1} onClick={() => setPage(page - 1)}>
                &lt;
            </button>

            Add clickable numbers for the offset here
            {[...Array(10).keys()].slice(0, 3).map(i =>
                <button key={i} onClick={() => setPage(i + 1)}>
                    {i + 1}
                </button>
            )}
            <button onClick={() => setPage(page + 1)}>
                &gt;
            </button>

*/}

export default Leads
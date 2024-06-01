import { createSlice, createAsyncThunk } from '@reduxjs/toolkit'
import axios from 'axios'



export const getTrips = createAsyncThunk('/leads/content', async (page) => {
    const offset = (page - 1) * 10
	const response = await axios.get(`/api/trip/get-list-trip?offset=${offset}`, {})
	return response.data;
})

export const leadsSlice = createSlice({
    name: 'leads',
    initialState: {
        isLoading: false,
        leads : []
    },
    reducers: {

        addNewLead: (state, action) => {
            let {newLeadObj} = action.payload
            state.leads = [...state.leads, newLeadObj]
        },

        deleteLead: (state, action) => {
            let {index} = action.payload
            state.leads.splice(index, 1)
        }
    },

    extraReducers: {
		[getTrips.pending]: state => {
			state.isLoading = true
		},
		[getTrips.fulfilled]: (state, action) => {
			state.leads = action.payload.data
			state.isLoading = false
		},
		[getTrips.rejected]: state => {
			state.isLoading = false
		},
    }
})

export const { addNewLead, deleteLead } = leadsSlice.actions

export default leadsSlice.reducer
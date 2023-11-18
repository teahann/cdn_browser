import { createClient } from '@supabase/supabase-js'

const URL = import.meta.env.VITE_PSQL_URL
const KEY = import.meta.env.VITE_PSQL_KEY

console.log(URL, KEY)

const database = createClient(URL, KEY)

export default database

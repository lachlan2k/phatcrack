/* Do not change, this code is generated from Golang structs */

export interface AdminAgentCreateRequestDTO {
  name: string
}
export interface AdminAgentCreateResponseDTO {
  name: string
  id: string
  key: string
}
export interface AdminUserCreateRequestDTO {
  username: string
  password: string
  roles: string[]
}
export interface AdminUserCreateResponseDTO {
  username: string
  id: string
  roles: string[]
}
export interface AdminIsSetupCompleteResponseDTO {
  is_complete: boolean
}
export interface AdminConfigResponseDTO {
  is_setup_complete: boolean
  is_mfa_required: boolean
  require_password_change_on_first_login: boolean
}
export interface AdminConfigRequestDTO {
  is_mfa_required: boolean
  require_password_change_on_first_login: boolean
}
export interface HashcatParams {
  attack_mode: number
  hash_type: number
  mask: string
  mask_increment: boolean
  mask_increment_min: number
  mask_increment_max: number
  mask_custom_charsets: string[]
  wordlist_filenames: string[]
  rules_filenames: string[]
  additional_args: string[]
  optimized_kernels: boolean
  slow_candidates: boolean
}
export interface AttackDTO {
  id: string
  hashlist_id: string
  hashcat_params: HashcatParams
}
export interface JobCrackedHashDTO {
  hash: string
  plaintext_hex: string
}
export interface HashcatStatusDevice {
  device_id: number
  device_name: string
  device_type: string
  speed: number
  util: number
  temp: number
}
export interface HashcatStatusGuess {
  guess_base: string
  guess_base_count: number
  guess_base_offset: number
  guess_base_percent: number
  guess_mod: string
  guess_mod_count: number
  guess_mod_offset: number
  guess_mod_percent: number
  guess_mode: number
}
export interface HashcatStatus {
  original_line: string
  time: Time
  session: string
  guess: HashcatStatusGuess
  status: number
  target: string
  progress: number[]
  restore_point: number
  recovered_hashes: number[]
  recovered_salts: number[]
  rejected: number
  devices: HashcatStatusDevice[]
}
export interface JobRuntimeOutputLineDTO {
  stream: string
  line: string
}
export interface Time {}
export interface JobRuntimeDataDTO {
  start_request_time: Time
  started_time: Time
  stopped_time: Time
  status: string
  stop_reason: string
  error_string: string
  output_lines: JobRuntimeOutputLineDTO[]
  status_updates: HashcatStatus[]
  cracked_hashes: JobCrackedHashDTO[]
}
export interface JobDTO {
  id: string
  hashlist_version: number
  attack_id: string
  hashcat_params: HashcatParams
  target_hashes: string[]
  hash_type: number
  runtime_data: JobRuntimeDataDTO
  assigned_agent_id: string
}
export interface AttackWithJobsDTO {
  id: string
  hashlist_id: string
  hashcat_params: HashcatParams
  jobs: JobDTO[]
}
export interface AttackWithJobsMultipleDTO {
  attacks: AttackWithJobsDTO[]
}
export interface AttackMultipleDTO {
  attacks: AttackDTO[]
}
export interface AttackCreateRequestDTO {
  hashlist_id: string
  hashcat_params: HashcatParams
  start_immediately: boolean
  name: string
  description: string
}
export interface AttackStartResponseDTO {
  new_job_id: string[]
}
export interface AuthLoginRequestDTO {
  username: string
  password: string
}
export interface AuthCurrentUserDTO {
  id: string
  username: string
  roles: string[]
}
export interface AuthLoginResponseDTO {
  user: AuthCurrentUserDTO
}
export interface AuthWhoamiResponseDTO {
  user: AuthCurrentUserDTO
}
export interface AuthRefreshResponseDTO {
  user: AuthCurrentUserDTO
}
export interface HashType {
  id: number
  name: string
  category: string
  slow_hash: boolean
  password_len_min: number
  password_len_max: number
  is_salted: boolean
  kernel_types: string[]
  example_hash_format: string
  example_hash: string
  example_pass: string
  benchmark_mask: string
  benchmark_charset1: string
  autodetect_enabled: boolean
  self_test_enabled: boolean
  potfile_enabled: boolean
  custom_plugin: boolean
  plaintext_encoding: string[]
}
export interface HashTypesDTO {
  hashtypes: { [key: number]: HashType }
}
export interface DetectHashTypeRequestDTO {
  test_hash: string
  has_username: boolean
}
export interface DetectHashTypeResponseDTO {
  possible_types: number[]
}
export interface VerifyHashesRequestDTO {
  hashes: string[]
  hash_type: number
  has_usernames: boolean
}
export interface VerifyHashesResponseDTO {
  valid: boolean
}
export interface NormalizeHashesResponseDTO {
  valid: boolean
  normalized_hashes: string[]
}
export interface HashlistCreateRequestDTO {
  project_id: string
  name: string
  hash_type: number
  input_hashes: string[]
  has_usernames: boolean
}
export interface HashlistCreateResponseDTO {
  id: string
}
export interface HashlistHashDTO {
  input_hash: string
  normalized_hash: string
  is_cracked: boolean
  plaintext_hex: string
}
export interface HashlistDTO {
  id: string
  name: string
  time_created: number
  hash_type: number
  hashes: HashlistHashDTO[]
  version: number
}
export interface HashlistResponseMultipleDTO {
  hashlists: HashlistDTO[]
}
export interface JobCreateRequestDTO {
  hashcat_params: HashcatParams
  hashes: string[]
  start_immediately: boolean
  name: string
  description: string
}
export interface JobCreateResponseDTO {
  id: string
}
export interface JobStartResponseDTO {
  agent_id: string
}

export interface JobSimpleDTO {
  id: string
  hashlist_version: number
  attack_id: string
  hash_type: number
  assigned_agent_id: string
}
export interface JobMultipleDTO {
  jobs: JobSimpleDTO[]
}
export interface WordlistCreateDTO {
  name: string
  description: string
  filename: string
  size: number
  lines: number
}
export interface RuleFileCreateDTO {
  name: string
  description: string
  filename: string
  size: number
  lines: number
}
export interface WordlistDTO {
  id: string
  name: string
  description: string
  filename_on_disk: string
  size_in_bytes: number
  lines: number
}
export interface RuleFileDTO {
  id: string
  name: string
  description: string
  filename_on_disk: string
  size_in_bytes: number
  lines: number
}
export interface GetAllWordlistsDTO {
  wordlists: WordlistDTO[]
}
export interface GetAllRuleFilesDTO {
  rulefiles: RuleFileDTO[]
}
export interface ProjectCreateDTO {
  name: string
  description: string
}
export interface ProjectDTO {
  id: string
  time_created: number
  name: string
  description: string
  owner_user_id: string
}
export interface ProjectResponseMultipleDTO {
  projects: ProjectDTO[]
}

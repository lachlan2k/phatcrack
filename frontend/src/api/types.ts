/* Do not change, this code is generated from Golang structs */

export interface AccountChangePasswordRequestDTO {
  new_password: string
  current_password: string
}
export interface AuthOIDCConfigDTO {
  client_id: string
  client_secret: string
  issuer_url: string
  redirect_url: string
  prompt: string
  automatic_user_creation: boolean
  username_claim: string
  role_field: string
  required_role: string
  scopes: string[]
}
export interface GeneralAuthConfigDTO {
  enabled_methods: string[]
  is_mfa_required: boolean
  require_password_change_on_first_login: boolean
}
export interface AuthConfigDTO {
  general?: GeneralAuthConfigDTO
  oidc?: AuthOIDCConfigDTO
}
export interface AgentConfigDTO {
  auto_sync_listfiles: boolean
  split_jobs_per_agent: number
}
export interface GeneralConfigDTO {
  is_maintenance_mode: boolean
  maximum_uploaded_file_size: number
  maximum_uploaded_file_line_scan_size: number
}
export interface AdminConfigRequestDTO {
  auth?: AuthConfigDTO
  agent?: AgentConfigDTO
  general?: GeneralConfigDTO
}
export interface AdminConfigResponseDTO {
  auth: AuthConfigDTO
  agent: AgentConfigDTO
  general: GeneralConfigDTO
}
export interface AdminAgentCreateRequestDTO {
  name: string
  ephemeral: boolean
}
export interface AdminAgentCreateResponseDTO {
  ephemeral: boolean
  name: string
  id: string
  key: string
}
export interface AdminAgentRegistrationKeyCreateRequestDTO {
  name: string
  ephemeral: boolean
}
export interface AdminAgentRegistrationKeyCreateResponseDTO {
  ephemeral: boolean
  name: string
  id: string
  key: string
}
export interface AdminUserCreateRequestDTO {
  username: string
  password: string
  gen_password: boolean
  lock_password: boolean
  roles: string[]
}
export interface AdminUserCreateResponseDTO {
  username: string
  id: string
  roles: string[]
  generated_password: string
}
export interface AdminUserUpdatePasswordRequestDTO {
  action: string
}
export interface AdminUserUpdatePasswordResponseDTO {
  generated_password: string
}
export interface AdminUserUpdateRequestDTO {
  username: string
  roles: string[]
}
export interface AdminServiceAccountCreateRequestDTO {
  username: string
  roles: string[]
}
export interface AdminServiceAccountCreateResponseDTO {
  username: string
  id: string
  roles: string[]
  api_key: string
}
export interface AdminGetUserDTO {
  id: string
  username: string
  roles: string[]
  is_password_locked: boolean
}
export interface AdminGetAllUsersResponseDTO {
  users: AdminGetUserDTO[]
}

export interface AdminAgentSetMaintanceRequestDTO {
  is_maintenance_mode: boolean
}
export interface HashcatStatusDevice {
  device_id: number
  device_name: string
  device_type: string
  speed: number
  util: number
  temp: number
}
export interface AgentFileDTO {
  name: string
  size: number
}
export interface AgentInfoDTO {
  status: string
  version: string
  last_checkin?: number
  available_listfiles?: AgentFileDTO[]
  active_job_ids?: string[]
}
export interface AgentDTO {
  id: string
  name: string
  is_maintenance_mode: boolean
  agent_info: AgentInfoDTO
  agent_devices: HashcatStatusDevice[]
}

export interface AgentGetAllResponseDTO {
  agents: AgentDTO[]
}
export interface AgentRegisterRequestDTO {
  name: string
}
export interface AgentRegisterResponseDTO {
  name: string
  id: string
  key: string
}
export interface HashcatParams {
  attack_mode: number
  hash_type: number
  mask: string
  mask_increment: boolean
  mask_increment_min: number
  mask_increment_max: number
  mask_sharded_charset: string
  mask_custom_charsets: string[]
  wordlist_filenames: string[]
  rules_filenames: string[]
  additional_args: string[]
  optimized_kernels: boolean
  slow_candidates: boolean
  skip: number
  limit: number
}
export interface AttackDTO {
  id: string
  hashlist_id: string
  hashcat_params: HashcatParams
  is_distributed: boolean
  progress_string: string
}
export interface AttackIDTreeDTO {
  project_id: string
  hashlist_id: string
  attack_id: string
}
export interface JobRuntimeSummaryDTO {
  hashrate: number
  estimated_time_remaining: number
  percent_complete: number
  started_time: number
  stopped_time: number
  cmd_line: string
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
  time_start: number
  estimated_stop: number
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
  cmd_line: string
  output_lines: JobRuntimeOutputLineDTO[]
  status_updates: HashcatStatus[]
}
export interface JobDTO {
  id: string
  hashlist_version: number
  attack_id: string
  hashcat_params: HashcatParams
  target_hashes: string[]
  hash_type: number
  runtime_data: JobRuntimeDataDTO
  runtime_summary: JobRuntimeSummaryDTO
  assigned_agent_id: string
}
export interface AttackWithJobsDTO {
  id: string
  hashlist_id: string
  hashcat_params: HashcatParams
  is_distributed: boolean
  progress_string: string
  jobs: JobDTO[]
}
export interface AttackWithJobsMultipleDTO {
  attacks: AttackWithJobsDTO[]
}
export interface AttackMultipleDTO {
  attacks: AttackDTO[]
}
export interface AttackIDTreeMultipleDTO {
  attacks: AttackIDTreeDTO[]
}
export interface AttackCreateRequestDTO {
  hashlist_id: string
  hashcat_params: HashcatParams
  is_distributed: boolean
}
export interface AttackStartResponseDTO {
  new_job_ids: string[]
  still_processing: boolean
}
export interface AuthLoginRequestDTO {
  username: string
  password: string
}
export interface AuthCurrentUserDTO {
  id: string
  username: string
  roles: string[]
  is_password_locked: boolean
}
export interface AuthLoginResponseDTO {
  user: AuthCurrentUserDTO
  is_awaiting_mfa: boolean
  requires_password_change: boolean
  requires_mfa_enrollment: boolean
}
export interface AuthWhoamiResponseDTO {
  user: AuthCurrentUserDTO
  is_awaiting_mfa: boolean
  requires_password_change: boolean
  requires_mfa_enrollment: boolean
}
export interface AuthRefreshResponseDTO {
  user: AuthCurrentUserDTO
  is_awaiting_mfa: boolean
  requires_password_change: boolean
  requires_mfa_enrollment: boolean
}
export interface AuthenticatorSelection {
  authenticatorAttachment?: string
  requireResidentKey?: boolean
  residentKey?: string
  userVerification?: string
}
export interface CredentialDescriptor {
  type: string
  id: number[]
  transports?: string[]
}
export interface CredentialParameter {
  type: string
  alg: number
}
export interface UserEntity {
  name: string
  displayName: string
  id: any
}
export interface RelyingPartyEntity {
  name: string
  id: string
}
export interface PublicKeyCredentialCreationOptions {
  rp: RelyingPartyEntity
  user: UserEntity
  challenge: number[]
  pubKeyCredParams?: CredentialParameter[]
  timeout?: number
  excludeCredentials?: CredentialDescriptor[]
  authenticatorSelection?: AuthenticatorSelection
  hints?: string[]
  attestation?: string
  attestationFormats?: string[]
  extensions?: { [key: string]: any }
}
export interface AuthWebAuthnStartEnrollmentResponseDTO {
  publicKey: PublicKeyCredentialCreationOptions
}
export interface PublicKeyCredentialRequestOptions {
  challenge: number[]
  timeout?: number
  rpId?: string
  allowCredentials?: CredentialDescriptor[]
  userVerification?: string
  hints?: string[]
  extensions?: { [key: string]: any }
}
export interface AuthWebAuthnStartChallengeResponseDTO {
  publicKey: PublicKeyCredentialRequestOptions
}
export interface AuthChangePasswordRequestDTO {
  old_password: string
  new_password: string
}
export interface PublicOIDCConfigDTO {
  prompt: string
}
export interface PublicAuthConfigDTO {
  enabled_methods: string[]
  oidc: PublicOIDCConfigDTO
}
export interface PublicGeneralConfigDTO {
  is_maintenance_mode: boolean
  maximum_uploaded_file_size: number
  maximum_uploaded_file_line_scan_size: number
}
export interface PublicConfigDTO {
  auth: PublicAuthConfigDTO
  general: PublicGeneralConfigDTO
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
export interface HashlistAppendRequestDTO {
  input_hashes: string[]
}
export interface HashlistAppendResponseDTO {
  num_new_hashes: number
  num_populated_from_potfile: number
}
export interface HashlistCreateResponseDTO {
  id: string
  num_populated_from_potfile: number
}
export interface HashlistHashDTO {
  id: string
  username: string
  input_hash: string
  normalized_hash: string
  is_cracked: boolean
  is_unexpected: boolean
  plaintext_hex: string
}
export interface HashlistDTO {
  id: string
  project_id: string
  name: string
  time_created: number
  hash_type: number
  hashes: HashlistHashDTO[]
  version: number
  has_usernames: boolean
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
export interface RunningJobForUserDTO {
  project_id: string
  hashlist_id: string
  attack_id: string
  job_id: string
}
export interface RunningJobsForUserResponseDTO {
  jobs: RunningJobForUserDTO[]
}
export interface RunningJobCountForUserDTO {
  username: string
  job_count: number
}
export interface RunningJobCountPerUsersDTO {
  result: RunningJobCountForUserDTO[]
}
export interface ListfileDTO {
  id: string
  file_type: string
  name: string
  size_in_bytes: number
  lines: number
  available_for_use: boolean
  pending_delete: boolean
  created_by_user_id: string
  associated_project_id: string
}
export interface GetAllWordlistsDTO {
  wordlists: ListfileDTO[]
}
export interface GetAllRuleFilesDTO {
  rulefiles: ListfileDTO[]
}
export interface GetAllListfilesDTO {
  listfiles: ListfileDTO[]
}
export interface ListfileUploadResponseDTO {
  listfile: ListfileDTO
}
export interface PotfileSearchRequestDTO {
  hashes: string[]
}
export interface PotfileSearchResultDTO {
  hash: string
  hash_type: number
  plaintext_hex: string
  found: boolean
}
export interface PotfileSearchResponseDTO {
  results: PotfileSearchResultDTO[]
}
export interface ProjectCreateRequestDTO {
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
export interface ProjectAddShareRequestDTO {
  user_id: string
}
export interface ProjectSharesDTO {
  user_ids: string[]
}
export interface UserDTO {
  id: string
  username: string
  roles: string[]
}
export interface UserMinimalDTO {
  id: string
  username: string
}
export interface UsersGetAllResponseDTO {
  users: UserMinimalDTO[]
}

/* Do not change, this code is generated from Golang structs */


export interface AdminAgentCreateRequestDTO {
    name: string;
}
export interface AdminAgentCreateResponseDTO {
    name: string;
    id: string;
    key: string;
}
export interface AdminUserCreateRequestDTO {
    username: string;
    password: string;
    role: string;
}
export interface AdminUserCreateResponseDTO {
    username: string;
    id: string;
    role: string;
}
export interface AuthLoginRequestDTO {
    username: string;
    password: string;
}
export interface AuthCurrentUserDTO {
    id: string;
    username: string;
    role: string;
}
export interface AuthLoginResponseDTO {
    user: AuthCurrentUserDTO;
}
export interface AuthWhoamiResponseDTO {
    user: AuthCurrentUserDTO;
}
export interface AuthRefreshResponseDTO {
    user: AuthCurrentUserDTO;
}
export interface HashcatParams {
    attack_mode: number;
    hash_type: number;
    mask: string;
    wordlist_filenames: string[];
    rules_filenames: string[];
    additional_args: string[];
    optimized_kernels: boolean;
    slow_candidates: boolean;
}
export interface JobCreateRequestDTO {
    hashcat_params: HashcatParams;
    hashes: string[];
    start_immediately: boolean;
    name: string;
    description: string;
}
export interface JobCreateResponseDTO {
    id: string;
}
export interface JobStartResponseDTO {
    agent_id: string;
}
export interface ListsWordlistCreateDTO {
    name: string;
    description: string;
    filename: string;
    size: number;
    lines: number;
}
export interface ListsRuleFileCreateDTO {
    name: string;
    description: string;
    filename: string;
    size: number;
    lines: number;
}
export interface ListsWordlistResponseDTO {
    name: string;
    description: string;
    filename: string;
    size: number;
    lines: number;
}
export interface ListsRuleFileResponseDTO {
    name: string;
    description: string;
    filename: string;
    size: number;
    lines: number;
}
export interface ListsGetAllWordlistsDTO {
    wordlists: ListsWordlistResponseDTO[];
}
export interface ListsGetAllRuleFilesDTO {
    rulefiles: ListsRuleFileResponseDTO[];
}
export interface ProjectCreateDTO {
    name: string;
    description: string;
}
export interface ProjectSimpleDetailsDTO {
    id: string;
    time_created: number;
    name: string;
    description: string;
}
export interface ProjectsFullDetailsDTO {
    id: string;
    time_created: number;
    name: string;
    description: string;
}
export interface ProjectResponseMultipleDTO {
    projects: ProjectSimpleDetailsDTO[];
}
export interface HashType {
    id: number;
    name: string;
    category: string;
    slow_hash: boolean;
    password_len_min: number;
    password_len_max: number;
    is_salted: boolean;
    kernel_types: string[];
    example_hash_format: string;
    example_hash: string;
    example_pass: string;
    benchmark_mask: string;
    benchmark_charset1: string;
    autodetect_enabled: boolean;
    self_test_enabled: boolean;
    potfile_enabled: boolean;
    custom_plugin: boolean;
    plaintext_encoding: string[];
}
export interface HashTypesDTO {
    hashtypes: {[key: number]: HashType};
}
export interface DetectHashTypeRequestDTO {
    test_hash: string;
}
export interface DetectHashTypeResponseDTO {
    possible_types: number[];
}
// Copyright 2019 Intel Corporation
// SPDX-License-Identifier: BSD-3-Clause
package deps

const compmappingJSON = `
[
    {"componentId":"intel_advisor", "dir":"advisor"},
    {"componentId":"collective_communications_library", "dir":"ccl"},
    {"componentId":"intel_cluster_checker", "dir":"clck"},
    {"componentId":"intel_data_analytics_library", "dir":"daal"},
    {"componentId":"1a_debugger", "dir":"debugger"},
    {"componentId":"dev-utilities", "dir":"dev-utilities"},
    {"componentId":"deep_neural_network", "dir":"oneDNN"},
    {"componentId":"dpcpp_compatibility_tool", "dir":"dpcpp-ct"},
    {"componentId":"eclipse-iot-plugins", "dir":"eclipse-iot-plugins"},
    {"componentId":"intel_embree", "dir":"embree"},
    {"componentId":"gstreamer_plugins", "dir":"gvaplugins"},
    {"componentId":"intel_inspector", "dir":"inspector"},
    {"componentId":"eclipse_based_ide", "dir":"intel-eclipse-ide"},
    {"componentId":"python", "dir":"intelpython"},
    {"componentId":"intel_integrated_performance_primitives", "dir":"ipp"},
    {"componentId":"intel_trace_analyzer_and_collector", "dir":"itac"},
    {"componentId":"intel_math_kernel_library", "dir":"mkl"},
    {"componentId":"intel_mpi_library", "dir":"mpi"},
    {"componentId":"open_image_denoise", "dir":"oidn"},
    {"componentId":"openvino", "dir":"openvino"},
    {"componentId":"open_volume_kernel_library", "dir":"openvkl"},
    {"componentId":"intel_ospray", "dir":"ospray"},
    {"componentId":"pstl", "dir":"pstl"},
    {"componentId":"intel_pytorch", "dir":"pytorch"},
    {"componentId":"intel_socwatch", "dir":"socwatch"},
    {"componentId":"intel_system_debugger", "dir":"system_debugger"},
    {"componentId":"intel_threading_building_blocks", "dir":"tbb"},
    {"componentId":"intel_tensorflow", "dir":"tensorflow"},
    {"componentId":"vpl", "dir":"vpl"},
    {"componentId":"intel_vtune_amplifier", "dir":"vtune"},


    {"componentId":"intel_cpp_compiler", "dir":"icc"},
    {"componentId":"openmp_fortran",     "dir":"fortran"},
    {"componentId":"dppcpp_compiler",    "dir":"dpcpp"},
    {"componentId":"dpcpp_library",      "dir":"?"},
    
    
    {"componentId":"intel_iot_connect_upm_mraa_cloud_connectors", "dir":"?"},
    {"componentId":"linux_kernel_build_tools", "dir":"?"},
    
    
    {"componentId":"gnu_project_debugger_gdb", "dir":"?"},

    {"componentId":"linux_iot_application_development_using_containerized_toolchains", "dir":"?"}
    

]
`

//if copy/pasting JSON from the WebConfigurator team, I recommend linting it first.
// https://jsonlint.com/
// they have lots of small problems (missing quotes, stray commas, missing commas) that
// Node.js forgives.
const sweetComponentsJSON = `
[
    
    { "suiteId": "HPCKit", "componentId": "intel_cpp_compiler", "primary": true },
    { "suiteId": "HPCKit", "componentId": "openmp_fortran", "primary": true },
    { "suiteId": "HPCKit", "componentId": "intel_mpi_library", "primary": true },
    { "suiteId": "HPCKit", "componentId": "intel_inspector", "primary": true },
    { "suiteId": "HPCKit", "componentId": "intel_trace_analyzer_and_collector", "primary": true },
    { "suiteId": "HPCKit", "componentId": "intel_cluster_checker", "primary": true },

    { "suiteId": "IOTKit", "componentId": "intel_cpp_compiler", "primary": true },
    { "suiteId": "IOTKit", "componentId": "intel_inspector", "primary": true },
    {
        "suiteId": "IOTKit",
        "componentId": "linux_iot_application_development_using_containerized_toolchains",
        "primary": true
    },
    { "suiteId": "IOTKit", "componentId": "eclipse_based_ide", "primary": true },
    { "suiteId": "IOTKit", "componentId": "linux_kernel_build_tools", "primary": true },
    { "suiteId": "IOTKit", "componentId": "intel_system_debugger", "primary": true },

    { "suiteId": "BringupKit", "componentId": "intel_socwatch", "primary": true },
    { "suiteId": "BringupKit", "componentId": "intel_system_debugger", "primary": true },


    { "suiteId": "DLDevKit", "componentId": "collective_communications_library", "primary": true },
    { "suiteId": "DLDevKit", "componentId": "deep_neural_network", "primary": true },

    { "suiteId": "AIKit", "componentId": "intel_tensorflow", "primary": true },
    { "suiteId": "AIKit", "componentId": "intel_pytorch", "primary": true },
    { "suiteId": "AIKit", "componentId": "python", "primary": true },

    { "suiteId": "oneAPIKit", "componentId": "dppcpp_compiler", "primary": true },
    { "suiteId": "oneAPIKit", "componentId": "dpcpp_compatibility_tool", "primary": true },
    { "suiteId": "oneAPIKit", "componentId": "dpcpp_library", "primary": true },
    { "suiteId": "oneAPIKit", "componentId": "1a_debugger", "primary": true },
    { "suiteId": "oneAPIKit", "componentId": "intel_math_kernel_library", "primary": true },
    { "suiteId": "oneAPIKit", "componentId": "intel_threading_building_blocks", "primary": true },
    { "suiteId": "oneAPIKit", "componentId": "intel_integrated_performance_primitives", "primary": true },
    { "suiteId": "oneAPIKit", "componentId": "intel_data_analytics_library", "primary": true },
    { "suiteId": "oneAPIKit", "componentId": "python", "primary": true },
    { "suiteId": "oneAPIKit", "componentId": "intel_advisor", "primary": true },
    { "suiteId": "oneAPIKit", "componentId": "deep_neural_network", "primary": true },
    { "suiteId": "oneAPIKit", "componentId": "collective_communications_library", "primary": true },
    { "suiteId": "oneAPIKit", "componentId": "vpl", "primary": true },

    { "suiteId": "RenderKit", "componentId": "intel_embree", "primary": true },
    { "suiteId": "RenderKit", "componentId": "intel_ospray", "primary": true },
    { "suiteId": "RenderKit", "componentId": "open_image_denoise", "primary": true },
    { "suiteId": "RenderKit", "componentId": "open_volume_kernel_library", "primary": true },

    { "suiteId": "VTuneProfiler", "componentId": "intel_vtune_amplifier", "primary": true }
]
`

const suitesJSON = `
[
    { "id": "HPCKit", "label": "Intel® oneAPI HPC Toolkit", "urlSlug": "hpc-kit", "baseToolkit": "dependency" },
    { "id": "IOTKit", "label": "Intel® oneAPI IoT Toolkit", "urlSlug": "iot-kit", "baseToolkit": "dependency" },
    { "id": "BringupKit", "label": "Intel® System Bring-Up Toolkit", "urlSlug": "bringup-kit", "baseToolkit": "recommended" },
    { "id": "DLDevKit", "label": "Intel® oneAPI DL Framework Developer Toolkit", "urlSlug": "dldev-kit", "baseToolkit": "dependency" },
    { "id": "AIKit", "label": "Intel® oneAPI AI Analytics Toolkit", "urlSlug": "ai-kit", "baseToolkit": "recommended" },
    { "id": "oneAPIKit", "label": "Intel® oneAPI Base Toolkit", "urlSlug": "oneapi-kit" },
    { "id": "RenderKit", "label": "Intel® oneAPI Rendering Toolkit", "urlSlug": "render-kit", "baseToolkit": "recommended" },
    { "id": "VTuneProfiler", "label": "Intel® VTune™ Profiler", "urlSlug": "vtune-profiler", "baseToolkit": "recommended" }
]
`

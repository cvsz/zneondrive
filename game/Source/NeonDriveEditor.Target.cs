using UnrealBuildTool;
using System.Collections.Generic;

public class NeonDriveEditorTarget : TargetRules
{
    public NeonDriveEditorTarget(TargetInfo Target) : base(Target)
    {
        Type = TargetType.Editor;
        DefaultBuildSettings = BuildSettingsVersion.V7;
        IncludeOrderVersion = EngineIncludeOrderVersion.Latest;
        ExtraModuleNames.Add("NeonDrive");
    }
}

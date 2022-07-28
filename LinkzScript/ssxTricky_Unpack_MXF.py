import struct


mxf_folder_path = "G:/Emulated/Xbox/Games/SSX Tricky Extracted.xiso/data/char/mdlxbx extracted/char/"
mxf_file_name = "mac_head" # without .mxf


modelHeaderList = [] # specify length by doing *18 or however many values there are


mxfFilePath = (mxf_folder_path+"/"+mxf_file_name+".mxf")
mxf = open(mxfFilePath, 'rb')

(fileUnk, modelCount, modelHeaderListOffset, modelListOffset) = struct.unpack('IHHI', mxf.read( 12)) # 12 bytes in file header


modelMegaList = []
#modelCount = 1 # CUSTOM MODEL COUNT

global modelHeaderOffsetShift
modelHeaderOffsetShift = 0

for m in range(modelCount): # Read Headers inside Model Header List

    mxf.seek(modelHeaderListOffset+modelHeaderOffsetShift,0) # SECTION: MODEL HEADER
    modelHeaderOffsetShift+=396

    modelName = mxf.read(16).decode('utf-8').strip('\x00')
    (modelOffset,            modelByteSize,             boneDataOffset,     boneDataOffset1,        
     materialListOffset,     boneDataOffset2,           ikDataListOffset,   morphSectionOffset,        
     skinningSectionOffset,  triStripSectionOffset,     unknown1,           vertexDataOffset,       
     unknown2,               unknownOffset1
    ) = struct.unpack('I'*14, mxf.read(4*14)) # unpack 14 32-bit integers

    mxf.seek(298, 1) # skip null values

    (unknownB,unknownB1,boneCount,morphHeaderCount,materialCount,ikDataCount,skinningHeaderCount,
     triStripGroupCount,unknownB2,vtxCount,unknownB3,unknownB4,unknownB5
    ) = struct.unpack('H'*13, mxf.read(2*13)) # unpack 14 16-bit integers


    currentModelOffset = modelListOffset + modelOffset

    materialListOffset    += currentModelOffset
    boneDataOffset        += currentModelOffset
    ikDataListOffset      += currentModelOffset
    morphSectionOffset    += currentModelOffset
    skinningSectionOffset += currentModelOffset
    triStripSectionOffset += currentModelOffset
    vertexDataOffset      += currentModelOffset



    #                                    SECTION: MATERIAL
    mxf.seek(materialListOffset, 0)

    materialList = []

    for mt in range(materialCount):
        texName                 = mxf.read(4).decode('utf-8').strip('\x00')
        texName1                = mxf.read(4).decode('utf-8').strip('\x00')
        texBumpName2            = mxf.read(4).decode('utf-8').strip('\x00')
        texG_Name3              = mxf.read(4).decode('utf-8').strip('\x00')
        texEnvName              = mxf.read(4).decode('utf-8').strip('\x00')
        (unkF0, unkF1, unkF2)   = struct.unpack('f'*3, mxf.read(4*3))
    # SECTION END: MATERIAL              #



    #                                    SECTION: BONE
    mxf.seek(boneDataOffset, 0)

    boneDataList = []

    for bn in range(boneCount):
        bone_name                       = mxf.read( 16).decode('utf-8').strip('\x00')
        (unk,boneParentID,unk1,boneID,) = struct.unpack('H'*4,   mxf.read(2*4)) # uint16
        BoneTranslation                 = struct.unpack('f'*3,   mxf.read(4*3))
        #mxf.seek(12*4, 1)
        (r, r1, r2, r3, r4, r5,
         u, u1, u2, u3, u4, u5)         = struct.unpack('f'*12, mxf.read(4*12)) # 12 single 32-bit floats
    # SECTION END: BONE



    #                                    SECTION: INVERSE KINEMATIC
    mxf.seek(ikDataListOffset, 0)

    ikDataList = []

    for ik in range(ikDataCount):
        ikLocation    = struct.unpack('f'*3, mxf.read(4*3)) # 3 floats
        (ikUnknown,)  = struct.unpack('I'  , mxf.read(  4)) # 1 integer
        ikDataList.append([ikLocation,ikUnknown])
    #print(f"\n{ikDataList}")
    # SECTION END: INVERSE KINEMATIC     #



    #                                    SECTION: MORPH
    mxf.seek(morphSectionOffset, 0)

    morphHeaderList = []
    morphContainer  = []

    global morphHeaderOffsetShift
    morphHeaderOffsetShift = 0

    for mh in range(morphHeaderCount):
        mxf.seek(morphSectionOffset+morphHeaderOffsetShift, 0)
        morphHeaderOffsetShift+=8

        (morphDataCount,
         morphDataListOffset) = struct.unpack('I'*2  , mxf.read(4*2)) # 1 integer

        morphHeaderList.append((morphDataCount, morphDataListOffset))

        mxf.seek(currentModelOffset+morphDataListOffset, 0)

        for md in range(morphDataCount):
            morphUnknown     = mxf.read(4)
            (morphUnknown1,) = struct.unpack('f', mxf.read(4))
            (morphUnknown2,) = struct.unpack('f', mxf.read(4))
            morphUnknown3    = mxf.read(4)
            #print(morphUnknown1)

            morphContainer.append([morphUnknown, morphUnknown1, morphUnknown2, morphUnknown3])
    # SECTION END: MORPH                 #



    #                                    SECTION: SKINNING
    mxf.seek(skinningSectionOffset, 0)

    skinningHeaderList    = []
    skinningDataContainer = []

    global skinningHeaderOffsetShift
    skinningHeaderOffsetShift = 0

    for sh in range(skinningHeaderCount):
        mxf.seek(skinningSectionOffset+skinningHeaderOffsetShift, 0)
        skinningHeaderOffsetShift+=12

        (skinningDataCount,
         skinningDataListOffset) = struct.unpack('I'*2, mxf.read(4*2))
        (skinningUnknown,
         skinningUnknown1)       = struct.unpack('H'*2, mxf.read(2*2))

        skinningHeaderList.append((skinningDataCount,skinningDataListOffset,skinningUnknown,skinningUnknown1))

        mxf.seek(currentModelOffset+skinningDataListOffset, 0)

        skinningTempList = []

        for sd in range(skinningDataCount):
            (boneWeightPercentage,
             boneID)       = struct.unpack('H'*2, mxf.read(2*2))
            skinningTempList.append((boneWeightPercentage, boneID))
        skinningDataContainer.append(skinningTempList)
    # SECTION END: SKINNING              #



    #                                    SECTION: TRISTRIP GROUP

    mxf.seek(triStripSectionOffset, 0)

    triStripHeaderList     = []
    triStripIndexContainer = []

    global triStripHeaderOffsetShift
    triStripHeaderOffsetShift = 0

    for th in range(triStripGroupCount):
        mxf.seek(triStripSectionOffset+triStripHeaderOffsetShift, 0)
        triStripHeaderOffsetShift+=16


        triStripheaderTuple = (
         indexListOffset, indexCount, materialIndex, materialIndex,  materialIndex,
         materialIndex,  materialIndex) = struct.unpack('IHHHHHH', mxf.read(16)) # 1 u32, 6 u16
        
        triStripHeaderList.append(triStripheaderTuple)


        mxf.seek(currentModelOffset+indexListOffset, 0)

        #for td in range(indexCount):


    print(triStripHeaderList)



#    print(f"""\n\n{modelName             = }
#{modelOffset           = }\n{modelByteSize         = }\n{boneDataOffset        = }\n{boneDataOffset1       = }
#{materialListOffset    = }\n{boneDataOffset2       = }\n{ikDataListOffset      = }\n{morphSectionOffset    = }
#{skinningSectionOffset = }\n{triStripSectionOffset = }\n{unknown1              = }\n{vertexDataOffset      = }
#{unknown2              = }\n{unknownOffset1        = }""")
#    print(f"""{unknownB            = }\n{unknownB1           = }\n{boneCount           = }\n{morphHeaderCount    = }
#{materialCount       = }\n{ikDataCount         = }\n{skinningHeaderCount = }\n{triStripGroupCount  = }
#{unknownB2           = }\n{vtxCount            = }\n{unknownB3           = }\n{unknownB4           = }
#{unknownB5           = }""")




    """

    if curmodelByteSize > 0:

        triStripHeaderList = []

        for trihdr in range(curTriStripHeaderCount):
            (triStripIndexListOffset            ,) = struct.unpack('I', mxf.read( 4))
            (triStripIndexCount                 ,) = struct.unpack('H', mxf.read( 2))
            mxf.seek(0xA, 1)
            #print(f"           {triStripIndexListOffset, triStripIndexCount}")
    
            triStripIndexListOffset += curModelOffset
    
            tmpB = [triStripIndexListOffset, triStripIndexCount]
            triStripHeaderList.append(tmpB)

        triStripsList = []
    
        for strips in range(curTriStripHeaderCount):
    
            curStripOffset  = triStripHeaderList[strips][0]
            curStripCount   = triStripHeaderList[strips][1]
    
            triStripIndexTuple = ()
    
            mxf.seek(curStripOffset,0)
            for index in range(curStripCount):
                (triStripIndex, ) = struct.unpack('H', mxf.read(2))
                triStripIndexTuple += (triStripIndex, ) 
            
            triStripsList.append(triStripIndexTuple)

    
        vtxStringList = []
        nrmStringList = []
        uvStringList  = []
    
        reOrient = False
    
        mxf.seek(curVtxOffset, 0)

        for v in range(curVtxCount):
            (vtxX, vtxY, vtxZ)          = struct.unpack('f'*3, mxf.read(4*3))
            (unkn,)                     = struct.unpack('f',   mxf.read(  4)) # unknown 
            (nrmX, nrmY, nrmZ)          = struct.unpack('f'*3, mxf.read(4*3))
            (unkn1,)                    = struct.unpack('I',   mxf.read(  4)) # unknown 
            (unkX, unkY, unkZ)          = struct.unpack('f'*3, mxf.read(4*3))
            (unkn2,)                    = struct.unpack('f',   mxf.read(  4)) # unknown 
            (uvMapU, uvMapV)            = struct.unpack('f'*2, mxf.read(4*2))
            (unkn3,)                    = struct.unpack('I',   mxf.read(  4)) # 0xFFFFFFFF
            (unkn4,)                    = struct.unpack('I',   mxf.read(  4)) # unknown 

            if reOrient:
                vtxStringList.append (f"\nv {round(  vtxX,8)} {round(       vtxZ,8)} {round(-vtxY,8)}")
                nrmStringList.append(f"\nvn {round(  nrmX,8)} {round(       nrmZ,8)} {round(-nrmY,8)}")
            else:
                vtxStringList.append (f"\nv {round(  vtxX,8)} {round(       vtxY,8)} {round(vtxZ,8)}")
                nrmStringList.append(f"\nvn {round(  nrmX,8)} {round(       nrmY,8)} {round(nrmZ,8)}")
            uvStringList.append (f"\nvt {round(uvMapU,8)} {round(-uvMapV+1.0,8)}")


        printInfo = True

        if printInfo:
            print(f"\nINFO:")
            print(f"\n  Model {mdl} | {hex(curModelOffset)} | {curModelName}:")
            print(f"\n          TriStrip Header Offset:{hex(curTriStripHeadOffset)}")
            print(f"            TriStrip Count:{curTriStripHeaderCount}\n")
            print(f"\n              Strip {strips} | Offset:{hex(curStripOffset)}  Count:{curStripCount}")
            print(f"\n          Vertices | {hex(curVtxOffset)} | Count:{curVtxCount}")
            print("\n                                                                       Model End\n\n")
    else:
        print(f"     Ignored: {mxf_file_name}_{mdl}_{curModelName} (Empty model)")
    """
    
mxf.close()

